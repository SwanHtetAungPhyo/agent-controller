package types

import "encoding/json"

type ClerkWebhookEvent struct {
	Data            json.RawMessage `json:"data"`
	EventAttributes EventAttributes `json:"event_attributes"`
	Object          string          `json:"object"`
	Timestamp       int64           `json:"timestamp"`
	Type            string          `json:"type"`
}

type EventAttributes struct {
	HTTPRequest HTTPRequest `json:"http_request"`
}

type HTTPRequest struct {
	ClientIP  string `json:"client_ip"`
	UserAgent string `json:"user_agent"`
}

type UserData struct {
	ID              string         `json:"id"`
	FirstName       string         `json:"first_name"`
	LastName        string         `json:"last_name"`
	EmailAddresses  []EmailAddress `json:"email_addresses"`
	ImageURL        string         `json:"image_url"`
	ProfileImageURL string         `json:"profile_image_url"`
	Username        string         `json:"username"`
	CreatedAt       int64          `json:"created_at"`
	UpdatedAt       int64          `json:"updated_at"`
}

type EmailAddress struct {
	ID           string `json:"id"`
	EmailAddress string `json:"email_address"`
}

type DeletedUserData struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
	Object  string `json:"object"`
}
