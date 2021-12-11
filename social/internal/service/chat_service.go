package service

import (
	"example.com/social/internal/domain"
	"github.com/godruoyi/go-snowflake"
)

type ChatService struct {
	Storage ChatStorage
}

func NewChatService(storage ChatStorage) *ChatService {
	return &ChatService{
		Storage: storage,
	}
}

func (chatService *ChatService) Publish(authorId, chatId int64, messageId string) error {
	message := &domain.Message{
		AuthorId: authorId,
		ChatId: chatId,
		Message: messageId,
		MessageId: snowflake.ID(),
	}
	return chatService.Storage.SaveMessage(message)
}