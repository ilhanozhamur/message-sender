package infrastructure

import (
	"context"
	"fmt"
	"github.com/ilhanozhamur/message-sender/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoDB struct {
	client     *mongo.Client
	collection *mongo.Collection
	state      string
}

type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Content   string             `bson:"content"`
	Recipient string             `bson:"recipient"`
	Sent      bool               `bson:"sent"`
	SentAt    *time.Time         `bson:"sentAt,omitempty"`
}

func NewMongoDB(uri, dbName, collectionName string) (*MongoDB, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to MongoDB")
	collection := client.Database(dbName).Collection(collectionName)

	return &MongoDB{
		client:     client,
		collection: collection,
		state:      "off",
	}, nil
}

func (db *MongoDB) SetState(state string) {
	db.state = state
}

func (db *MongoDB) GetState() string {
	return db.state
}

func (db *MongoDB) GetUnsentMessages(limit int) []*domain.Message {
	filter := bson.M{"sent": false}
	cursor, err := db.collection.Find(context.TODO(), filter, options.Find().SetLimit(int64(limit)))
	if err != nil {
		return nil
	}
	defer cursor.Close(context.TODO())

	var unsentMessages []*domain.Message
	for cursor.Next(context.TODO()) {
		var msg Message
		if err := cursor.Decode(&msg); err != nil {
			continue
		}
		unsentMessages = append(unsentMessages, &domain.Message{
			ID:        msg.ID,
			Content:   msg.Content,
			Recipient: msg.Recipient,
			Sent:      msg.Sent,
			SentAt:    msg.SentAt,
		})
	}
	return unsentMessages
}

func (db *MongoDB) MarkMessageAsSent(id primitive.ObjectID) error {
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"sent":   true,
			"sentAt": now,
		},
	}
	_, err := db.collection.UpdateOne(context.TODO(), bson.M{"_id": id}, update)
	return err
}

func (db *MongoDB) GetSentMessages() ([]domain.Message, error) {
	cursor, err := db.collection.Find(context.TODO(), bson.M{"sent": true})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var sentMessages []domain.Message
	for cursor.Next(context.TODO()) {
		var msg Message
		if err := cursor.Decode(&msg); err != nil {
			continue
		}
		sentMessages = append(sentMessages, domain.Message{
			ID:        msg.ID,
			Content:   msg.Content,
			Recipient: msg.Recipient,
			Sent:      msg.Sent,
			SentAt:    msg.SentAt,
		})
	}
	return sentMessages, nil
}
