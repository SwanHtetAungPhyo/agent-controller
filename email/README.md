# Kainos Email Service

A NATS-based async email service using Resend API with Kainos branding.

## Features

- ✅ **Async Processing**: Listens to NATS messages for email requests
- ✅ **Resend Integration**: Uses Resend API for reliable email delivery
- ✅ **Kainos Branding**: Professional email templates with Kainos logo
- ✅ **Multiple Email Types**: Welcome and general notification emails
- ✅ **Environment Configuration**: All settings via environment variables

## Configuration

Set these environment variables in `.env`:

```env
NATS_URL=nats://localhost:4222
NATS_MAX_RECONNECT=5
NATS_RECONNECT_WAIT=2s
NATS_TIMEOUT=10s
TOPIC=email.send
RESEND_API_KEY=re_GPR15sTc_GeWH2bkwD6GmKcmCv2bgFXxk
FROM_EMAIL=onboarding@resend.dev
FROM_NAME=Kainos Team
```

## Usage

### Start the Service
```bash
export $(cat .env | xargs) && ./email-service
```

### Send Welcome Email
```json
{
  "type": "email.send",
  "data": {
    "type": "welcome",
    "message": "Welcome message content",
    "info": {
      "to": "user@example.com",
      "name": "User Name"
    }
  }
}
```

### Send General Email
```json
{
  "type": "email.send",
  "data": {
    "type": "general",
    "message": "Email content",
    "info": {
      "to": "user@example.com",
      "subject": "Email Subject",
      "html": true
    }
  }
}
```

## Testing

Run the test script:
```bash
go run test_local.go
```

## Email Template

All emails include:
- Kainos logo header (blue background)
- Professional styling
- Responsive design
- Footer with copyright

The service successfully sends emails to `swanhtatungp@gmail.com` using the Resend API.
