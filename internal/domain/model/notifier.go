package model

type TelegramService interface {
	SendMessageToUserByTelegramID(userID TelegramID, text string) error
	SendMessageToUserByUserID(userID UserID, text string) error
}
