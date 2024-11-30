package db

import "context"

type Message struct {
	Version       string `json:"_Version" xml:"Version,attr"`
	Message       string `json:"_Message" xml:"Message,attr"`
	MessageTypeId string `json:"_MessageTypeId" xml:"MessageTypeId,attr"`
	MessageId     string `json:"_Message_id" xml:"MessageId,attr"`
	RecurringType string `json:"recurringType" xml:"recurringType"`
	Currency      string `json:"currency" xml:"currency"`
	PrvId         string `json:"prv_id" xml:"prv_id"`
}

type DatabaseService interface {
	InsertMessage(ctx context.Context, msg string) error
	InsertMessages(ctx context.Context, msgs []string) error

	InsertMessageJson(ctx context.Context, msg string) error
	InsertMessagesJson(ctx context.Context, msgs []string) error
}
