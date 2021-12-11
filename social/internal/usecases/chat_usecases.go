package usecases

type PublishMessageUseCase interface {
	Publish(authorId, chatId int64, messageId string) error
}