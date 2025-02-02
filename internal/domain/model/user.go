package model

import "time"

type UserID int64
type TelegramID int64

type User struct {
	ID         UserID     `json:"user_id"`
	TelegramID TelegramID `json:"telegram_id"`
	Username   string     `json:"username"`
	CreatedAt  time.Time  `json:"created_at"`
}

func (u *User) CanReceiveNotifications() bool {
	return true
}
