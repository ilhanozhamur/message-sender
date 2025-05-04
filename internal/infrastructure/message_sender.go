package infrastructure

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
)

type MessageSender struct {
	db          *MongoDB
	redis       *RedisClient
	ticker      *time.Ticker
	quit        chan struct{}
	stopChannel chan struct{}
}

func NewMessageSender(db *MongoDB, redis *RedisClient) *MessageSender {
	return &MessageSender{
		db:          db,
		redis:       redis,
		ticker:      time.NewTicker(2 * time.Minute),
		quit:        make(chan struct{}),
		stopChannel: make(chan struct{}),
	}
}

func (ms *MessageSender) Start() {
	for {
		select {
		case <-ms.ticker.C:
			ms.sendUnsentMessages()
		case <-ms.quit:
			ms.ticker.Stop()
			return
		}
	}
}

func (ms *MessageSender) Stop() {
	select {
	case <-ms.stopChannel:
		return
	default:
		close(ms.stopChannel)
	}
}

func (ms *MessageSender) sendUnsentMessages() {
	messages := ms.db.GetUnsentMessages(10)
	for _, message := range messages {
		infraMessage := &Message{
			ID:        message.ID,
			Content:   message.Content,
			Recipient: message.Recipient,
			Sent:      message.Sent,
			SentAt:    message.SentAt,
		}
		if err := ms.sendMessage(infraMessage); err != nil {
			log.Println("Failed to send message:", err)
		}
	}
}

func (ms *MessageSender) sendMessage(message *Message) error {
	payload := map[string]string{
		"to":      message.Recipient,
		"content": message.Content,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return errors.New("failed to marshal payload: " + err.Error())
	}

	req, err := http.NewRequest("POST", "https://webhook.site/cf524fbe-e4c1-40a5-a96d-1675692a3be7", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return errors.New("failed to create HTTP request: " + err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-ins-auth-key", "INS.me1x9uMcyYGlhKKQVPoc.bO3j9aZwRTOcA2Ywb")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.New("error sending HTTP request: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return errors.New("unexpected response status: " + resp.Status)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return errors.New("failed to decode response body: " + err.Error())
	}

	messageID, ok := response["messageId"].(string)
	if !ok {
		return errors.New("response does not contain a valid messageId")
	}

	ms.db.MarkMessageAsSent(message.ID)
	if ms.redis != nil {

		ms.redis.CacheMessage(messageID, time.Now())
		keys, err := ms.redis.GetKeys("*")
		if err != nil {
			log.Println("Error fetching keys:", err)
		} else {
			log.Println("Redis keys:", keys)
		}
	}

	return nil
}
