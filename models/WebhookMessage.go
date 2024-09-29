package models

type WebhookMessage struct {
	Event              string `json:"event"`
	APIKey             string `json:"apikey"`
	PhoneNumber        string `json:"phone_number"`
	ContactPhoneNumber string `json:"contact_phone_number"`
	ContactName        string `json:"contact_name"`
	MessageBody        string `json:"message_body"`
	MessageCustomID    string `json:"message_custom_id"`
	EventTime          string `json:"event_time"`
}