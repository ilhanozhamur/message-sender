package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Content   string             `bson:"content"`
	Recipient string             `bson:"recipient"`
	Sent      bool               `bson:"sent"`
	SentAt    *time.Time         `bson:"sentAt,omitempty"`
}

type MessageRepository interface {
	GetUnsentMessages(limit int) []*Message
	MarkMessageAsSent(id primitive.ObjectID) error
	GetSentMessages() ([]Message, error)
	SetState(state string)
	GetState() string
}

type MessageSender interface {
	Start()
	Stop()
}
