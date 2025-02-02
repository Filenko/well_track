package model

import (
	"errors"
	"time"
)

type AnswerID int64

var ErrInvalidRating = errors.New("rating must be between 1 and 5")

type Answer struct {
	ID        AnswerID  `json:"answer_id"`
	UserID    UserID    `json:"user_id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created-at"`
}

func (a *Answer) Validate() error {
	if a.Rating < 1 || a.Rating > 5 {
		return ErrInvalidRating
	}
	return nil
}
