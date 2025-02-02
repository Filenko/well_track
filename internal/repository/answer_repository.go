package repository

import (
	"well_track/internal/domain/model"
)

type AnswerRepository interface {
	Create(answer *model.Answer) error
	//TODO: implement GetByUserID
	//GetByUserID(userID model.UserID) ([]model.Answer, error)
}
