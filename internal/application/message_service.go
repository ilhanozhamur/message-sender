package application

import (
	"github.com/ilhanozhamur/message-sender/internal/domain"
)

type MessageService struct {
	repo   domain.MessageRepository
	sender domain.MessageSender
}

func NewMessageService(repo domain.MessageRepository, sender domain.MessageSender) *MessageService {
	return &MessageService{
		repo:   repo,
		sender: sender,
	}
}

func (ms *MessageService) StartMessageSending() {
	ms.repo.SetState("on")
	ms.sender.Start()
}

func (ms *MessageService) StopMessageSending() {
	ms.repo.SetState("off")
	ms.sender.Stop()
}

func (ms *MessageService) GetMessageSendingState() string {
	return ms.repo.GetState()
}

func (ms *MessageService) GetSentMessages() ([]domain.Message, error) {
	return ms.repo.GetSentMessages()
}
