package repository

import (
	"well_track/internal/domain/model"
)

type UserRepository interface {
	Create(user *model.User) error
	GetByID(userID model.UserID) (*model.User, error)
	GetByTelegramID(tgID model.TelegramID) (*model.User, error)
	//TODO: implement Update
	//Update(user *model.User) error
}
