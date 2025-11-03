package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nats-io/nats.go"
)

type Payload struct {
	Type    string                 `json:"type"`
	Message string                 `json:"message"`
	Info    map[string]interface{} `json:"info"`
}

type Config struct {
	ResendAPIKey string
	FromEmail    string
	FromName     string
}

type EmailService struct {
	config Config
	nc     *nats.Conn
	js     nats.JetStreamContext
	client *http.Client
}

type ResendEmail struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
}

type ResendResponse struct {
	ID string `json:"id"`
}

func NewEmailService(cfg Config, nc *nats.Conn) (*EmailService, error) {
	js, err := nc.JetStream()
	if err != nil {
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	return &EmailService{
		config: cfg,
		nc:     nc,
		js:     js,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

func (es *EmailService) SendEmail(to []string, subject, htmlBody string) error {
	if len(to) == 0 {
		return fmt.Errorf("no recipients specified")
	}

	fromAddress := fmt.Sprintf("%s <%s>", es.config.FromName, es.config.FromEmail)

	email := ResendEmail{
		From:    fromAddress,
		To:      to,
		Subject: subject,
		HTML:    htmlBody,
	}

	return es.sendWithResend(email)
}

func (es *EmailService) sendWithResend(email ResendEmail) error {
	jsonData, err := json.Marshal(email)
	if err != nil {
		return fmt.Errorf("failed to marshal email: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+es.config.ResendAPIKey)

	resp, err := es.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var buf bytes.Buffer
		buf.ReadFrom(resp.Body)
		return fmt.Errorf("resend API error: %d - %s", resp.StatusCode, buf.String())
	}

	var resendResp ResendResponse
	if err := json.NewDecoder(resp.Body).Decode(&resendResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	log.Printf("Email sent successfully with ID: %s", resendResp.ID)
	return nil
}

func (es *EmailService) StartListener(ctx context.Context, subject string) error {
	if es.nc == nil || !es.nc.IsConnected() {
		return fmt.Errorf("NATS connection is not available")
	}

	if es.js != nil {
		if err := es.setupJetStream(); err != nil {
			log.Printf("JetStream setup failed, falling back to regular NATS: %v", err)
			return es.startRegularNATSListener(subject)
		}

		_, err := es.js.Subscribe(subject, func(msg *nats.Msg) {
			var event struct {
				Type string                 `json:"type"`
				Data map[string]interface{} `json:"data"`
			}
			if err := json.Unmarshal(msg.Data, &event); err != nil {
				log.Printf("Failed to unmarshal event: %v", err)
				msg.Nak()
				return
			}

			payload := &Payload{
				Type:    event.Data["type"].(string),
				Message: event.Data["message"].(string),
				Info:    event.Data["info"].(map[string]interface{}),
			}

			es.processEmailPayload(payload)
			msg.Ack()
		}, nats.Durable("email-consumer"), nats.ManualAck())

		if err != nil {
			log.Printf("JetStream subscription failed, falling back to regular NATS: %v", err)
			return es.startRegularNATSListener(subject)
		}

		log.Printf("Email service listening on JetStream subject: %s with stream: event_email", subject)
		return nil
	}

	return es.startRegularNATSListener(subject)
}

func (es *EmailService) startRegularNATSListener(subject string) error {
	_, err := es.nc.Subscribe(subject, func(msg *nats.Msg) {
		var event struct {
			Type string                 `json:"type"`
			Data map[string]interface{} `json:"data"`
		}
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("Failed to unmarshal event: %v", err)
			return
		}

		payload := &Payload{
			Type:    event.Data["type"].(string),
			Message: event.Data["message"].(string),
			Info:    event.Data["info"].(map[string]interface{}),
		}

		go es.processEmailPayload(payload)
	})

	if err != nil {
		return fmt.Errorf("failed to subscribe to regular NATS subject: %w", err)
	}

	log.Printf("Email service listening on regular NATS subject: %s", subject)
	return nil
}

func (es *EmailService) setupJetStream() error {
	streamName := "event_email"

	stream, err := es.js.StreamInfo(streamName)
	if err != nil {
		if err == nats.ErrStreamNotFound {
			log.Printf("Creating JetStream stream: %s", streamName)
			_, err = es.js.AddStream(&nats.StreamConfig{
				Name:     streamName,
				Subjects: []string{"email.>"},
				Storage:  nats.FileStorage,
				MaxAge:   24 * time.Hour,
				Replicas: 1,
			})
			if err != nil {
				return fmt.Errorf("failed to create stream: %w", err)
			}
			log.Printf("JetStream stream '%s' created successfully", streamName)
		} else {
			return fmt.Errorf("failed to get stream info: %w", err)
		}
	} else {
		log.Printf("JetStream stream '%s' already exists with %d messages", streamName, stream.State.Msgs)
	}

	return nil
}

func (es *EmailService) processEmailPayload(payload *Payload) {
	switch payload.Type {
	case "welcome":
		es.sendWelcomeEmail(payload)
	case "general":
		es.sendGeneralEmail(payload)
	default:
		log.Printf("Unknown email type: %s", payload.Type)
	}
}

func (es *EmailService) sendWelcomeEmail(payload *Payload) {
	to, ok := payload.Info["to"].(string)
	if !ok {
		log.Printf("Missing 'to' field in welcome email payload")
		return
	}

	name, _ := payload.Info["name"].(string)
	if name == "" {
		name = "User"
	}

	subject := "Welcome to Kainos!"
	htmlBody := es.buildEmailTemplate("Welcome "+name+"!", payload.Message, `
		<div style="text-align: center; margin: 20px 0;">
			<p style="font-size: 18px; color: #333;">Thank you for joining us at Kainos!</p>
			<p style="color: #666;">We're excited to have you on board and look forward to working with you.</p>
		</div>
	`)

	if err := es.SendEmail([]string{to}, subject, htmlBody); err != nil {
		log.Printf("Failed to send welcome email: %v", err)
	} else {
		log.Printf("Welcome email sent to: %s", to)
	}
}

func (es *EmailService) sendGeneralEmail(payload *Payload) {
	to, ok := payload.Info["to"].(string)
	if !ok {
		log.Printf("Missing 'to' field in general email payload")
		return
	}

	subject, _ := payload.Info["subject"].(string)
	if subject == "" {
		subject = "Kainos Notification"
	}

	htmlBody := es.buildEmailTemplate("Notification", payload.Message, "")

	if err := es.SendEmail([]string{to}, subject, htmlBody); err != nil {
		log.Printf("Failed to send general email: %v", err)
	} else {
		log.Printf("General email sent to: %s", to)
	}
}

// buildEmailTemplate creates a beautiful HTML email template with cyan theme
func (es *EmailService) buildEmailTemplate(title, message, additionalContent string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            margin: 0;
            padding: 0;
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            min-height: 100vh;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            background-color: #ffffff;
            box-shadow: 0 10px 30px rgba(0,0,0,0.1);
            border-radius: 10px;
            overflow: hidden;
        }
        .header {
            background: linear-gradient(135deg, #00d4ff 0%%, #0099cc 100%%);
            padding: 30px 20px;
            text-align: center;
            position: relative;
        }
        .header::before {
            content: '';
            position: absolute;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background: url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><defs><pattern id="grain" width="100" height="100" patternUnits="userSpaceOnUse"><circle cx="25" cy="25" r="1" fill="rgba(255,255,255,0.1)"/><circle cx="75" cy="75" r="1" fill="rgba(255,255,255,0.1)"/><circle cx="50" cy="10" r="0.5" fill="rgba(255,255,255,0.1)"/></pattern></defs><rect width="100" height="100" fill="url(%%23grain)"/></svg>');
            opacity: 0.3;
        }
        .logo {
            color: #ffffff;
            font-size: 36px;
            font-weight: 700;
            margin: 0;
            text-shadow: 0 2px 4px rgba(0,0,0,0.2);
            letter-spacing: 2px;
            position: relative;
            z-index: 1;
        }
        .content {
            padding: 40px 30px;
            background: #ffffff;
        }
        .title {
            color: #2c3e50;
            font-size: 28px;
            margin-bottom: 25px;
            font-weight: 600;
            text-align: center;
        }
        .message {
            color: #34495e;
            font-size: 16px;
            line-height: 1.8;
            margin-bottom: 25px;
            text-align: center;
        }
        .cta-button {
            display: inline-block;
            background: linear-gradient(135deg, #00d4ff 0%%, #0099cc 100%%);
            color: white;
            padding: 15px 30px;
            text-decoration: none;
            border-radius: 25px;
            font-weight: 600;
            margin: 20px 0;
            box-shadow: 0 4px 15px rgba(0, 212, 255, 0.3);
            transition: all 0.3s ease;
        }
        .footer {
            background: linear-gradient(135deg, #f8f9fa 0%%, #e9ecef 100%%);
            padding: 25px 20px;
            text-align: center;
            border-top: 1px solid #dee2e6;
        }
        .footer-text {
            color: #6c757d;
            font-size: 14px;
            margin: 5px 0;
            line-height: 1.5;
        }
        .social-links {
            margin: 15px 0;
        }
        .social-link {
            display: inline-block;
            width: 40px;
            height: 40px;
            background: linear-gradient(135deg, #00d4ff 0%%, #0099cc 100%%);
            border-radius: 50%%;
            margin: 0 5px;
            line-height: 40px;
            color: white;
            text-decoration: none;
            font-weight: bold;
        }
        @media (max-width: 600px) {
            .container { margin: 10px; }
            .content { padding: 25px 20px; }
            .title { font-size: 24px; }
            .logo { font-size: 28px; }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1 class="logo">KAINOS</h1>
        </div>
        <div class="content">
            <h2 class="title">%s</h2>
            <div class="message">%s</div>
            %s
        </div>
        <div class="footer">
            <div class="social-links">
                <a href="#" class="social-link">K</a>
            </div>
            <p class="footer-text">© 2024 Kainos. All rights reserved.</p>
            <p class="footer-text">This email was sent from Kainos notification system.</p>
            <p class="footer-text">Building the future, one innovation at a time.</p>
        </div>
    </div>
</body>
</html>`, title, title, message, additionalContent)
}

func (es *EmailService) buildWelcomeEmailTemplate(name, message string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome to Kainos</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #00bcd4 0%%, #0097a7 100%%);
            padding: 20px;
        }
        .email-container {
            max-width: 600px;
            margin: 0 auto;
            background: #ffffff;
            border-radius: 16px;
            overflow: hidden;
            box-shadow: 0 20px 40px rgba(0, 188, 212, 0.15);
        }
        .header {
            background: linear-gradient(135deg, #00bcd4 0%%, #00acc1 50%%, #0097a7 100%%);
            padding: 40px 30px;
            text-align: center;
            position: relative;
        }
        .header::before {
            content: '';
            position: absolute;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background: url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><defs><pattern id="grain" width="100" height="100" patternUnits="userSpaceOnUse"><circle cx="25" cy="25" r="1" fill="rgba(255,255,255,0.1)"/><circle cx="75" cy="75" r="1" fill="rgba(255,255,255,0.1)"/><circle cx="50" cy="10" r="0.5" fill="rgba(255,255,255,0.05)"/></pattern></defs><rect width="100" height="100" fill="url(%%23grain)"/></svg>');
            opacity: 0.3;
        }
        .logo {
            color: #ffffff;
            font-size: 42px;
            font-weight: 700;
            margin: 0;
            text-shadow: 0 2px 4px rgba(0,0,0,0.1);
            letter-spacing: 2px;
            position: relative;
            z-index: 1;
        }
        .welcome-icon {
            font-size: 64px;
            margin-bottom: 20px;
            position: relative;
            z-index: 1;
        }
        .content {
            padding: 50px 40px;
            text-align: center;
        }
        .greeting {
            color: #00bcd4;
            font-size: 32px;
            font-weight: 600;
            margin-bottom: 20px;
            text-shadow: 0 1px 2px rgba(0,0,0,0.1);
        }
        .message {
            color: #37474f;
            font-size: 18px;
            line-height: 1.8;
            margin-bottom: 30px;
            max-width: 480px;
            margin-left: auto;
            margin-right: auto;
        }
        .highlight-box {
            background: linear-gradient(135deg, #e0f7fa 0%%, #b2ebf2 100%%);
            border-left: 4px solid #00bcd4;
            padding: 25px;
            margin: 30px 0;
            border-radius: 8px;
            text-align: left;
        }
        .highlight-title {
            color: #00695c;
            font-size: 20px;
            font-weight: 600;
            margin-bottom: 15px;
        }
        .highlight-text {
            color: #004d40;
            font-size: 16px;
            line-height: 1.6;
        }
        .cta-button {
            display: inline-block;
            background: linear-gradient(135deg, #00bcd4 0%%, #0097a7 100%%);
            color: #ffffff;
            text-decoration: none;
            padding: 16px 32px;
            border-radius: 50px;
            font-size: 18px;
            font-weight: 600;
            margin: 20px 0;
            box-shadow: 0 8px 20px rgba(0, 188, 212, 0.3);
            transition: all 0.3s ease;
        }
        .cta-button:hover {
            transform: translateY(-2px);
            box-shadow: 0 12px 25px rgba(0, 188, 212, 0.4);
        }
        .features {
            display: flex;
            justify-content: space-around;
            margin: 40px 0;
            flex-wrap: wrap;
        }
        .feature {
            text-align: center;
            flex: 1;
            min-width: 150px;
            margin: 10px;
        }
        .feature-icon {
            font-size: 48px;
            color: #00bcd4;
            margin-bottom: 15px;
        }
        .feature-title {
            color: #00695c;
            font-size: 16px;
            font-weight: 600;
            margin-bottom: 8px;
        }
        .feature-text {
            color: #546e7a;
            font-size: 14px;
            line-height: 1.4;
        }
        .footer {
            background: linear-gradient(135deg, #f0fdff 0%%, #e0f2f1 100%%);
            padding: 30px;
            text-align: center;
            border-top: 1px solid #b2dfdb;
        }
        .footer-text {
            color: #546e7a;
            font-size: 14px;
            margin: 5px 0;
            line-height: 1.5;
        }
        .social-links {
            margin: 20px 0;
        }
        .social-link {
            display: inline-block;
            margin: 0 10px;
            color: #00bcd4;
            text-decoration: none;
            font-weight: 500;
        }
        @media (max-width: 600px) {
            .content { padding: 30px 20px; }
            .greeting { font-size: 28px; }
            .message { font-size: 16px; }
            .features { flex-direction: column; }
        }
    </style>
</head>
<body>
    <div class="email-container">
        <div class="header">
            <div class="welcome-icon">🚀</div>
            <h1 class="logo">KAINOS</h1>
        </div>
        <div class="content">
            <h2 class="greeting">Welcome, %s! 👋</h2>
            <p class="message">%s</p>

            <div class="highlight-box">
                <div class="highlight-title">🎉 You're all set!</div>
                <div class="highlight-text">
                    Your account has been successfully created and you now have access to all Kainos features.
                    We're excited to have you join our community of innovators and creators.
                </div>
            </div>

            <a href="#" class="cta-button">Get Started Now</a>

            <div class="features">
                <div class="feature">
                    <div class="feature-icon">⚡</div>
                    <div class="feature-title">Lightning Fast</div>
                    <div class="feature-text">Experience blazing fast performance</div>
                </div>
                <div class="feature">
                    <div class="feature-icon">🔒</div>
                    <div class="feature-title">Secure</div>
                    <div class="feature-text">Your data is protected with enterprise-grade security</div>
                </div>
                <div class="feature">
                    <div class="feature-icon">🌟</div>
                    <div class="feature-title">Premium</div>
                    <div class="feature-text">Access to all premium features and support</div>
                </div>
            </div>
        </div>
        <div class="footer">
            <div class="social-links">
                <a href="#" class="social-link">📧 Support</a>
                <a href="#" class="social-link">📱 Mobile App</a>
                <a href="#" class="social-link">📚 Documentation</a>
            </div>
            <p class="footer-text">© 2024 Kainos. All rights reserved.</p>
            <p class="footer-text">This welcome email was sent because you just joined Kainos.</p>
            <p class="footer-text">If you have any questions, feel free to reach out to our support team.</p>
        </div>
    </div>
</body>
</html>`, name, message)
}

func (es *EmailService) buildGeneralEmailTemplate(title, message string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #e0f7fa 0%%, #b2ebf2 100%%);
            padding: 20px;
        }
        .email-container {
            max-width: 600px;
            margin: 0 auto;
            background: #ffffff;
            border-radius: 12px;
            overflow: hidden;
            box-shadow: 0 10px 30px rgba(0, 188, 212, 0.1);
        }
        .header {
            background: linear-gradient(135deg, #00bcd4 0%%, #0097a7 100%%);
            padding: 30px;
            text-align: center;
        }
        .logo {
            color: #ffffff;
            font-size: 28px;
            font-weight: 600;
            margin: 0;
            letter-spacing: 1px;
        }
        .content {
            padding: 40px 30px;
        }
        .title {
            color: #00695c;
            font-size: 24px;
            margin-bottom: 20px;
            text-align: center;
        }
        .message {
            color: #37474f;
            font-size: 16px;
            line-height: 1.7;
            margin-bottom: 20px;
            text-align: center;
        }
        .footer {
            background: #f0fdff;
            padding: 20px;
            text-align: center;
            border-top: 1px solid #b2dfdb;
        }
        .footer-text {
            color: #546e7a;
            font-size: 14px;
            margin: 5px 0;
        }
    </style>
</head>
<body>
    <div class="email-container">
        <div class="header">
            <h1 class="logo">KAINOS</h1>
        </div>
        <div class="content">
            <h2 class="title">%s</h2>
            <div class="message">%s</div>
        </div>
        <div class="footer">
            <p class="footer-text">© 2024 Kainos. All rights reserved.</p>
            <p class="footer-text">This email was sent from Kainos notification system.</p>
        </div>
    </div>
</body>
</html>`, title, title, message)
}

func (es *EmailService) Validate() error {
	if es.config.ResendAPIKey == "" {
		return fmt.Errorf("Resend API key is required")
	}
	if es.config.FromEmail == "" {
		return fmt.Errorf("from email is required")
	}
	return nil
}
