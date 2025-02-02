package usecase

import (
	"github.com/rs/zerolog"
	"time"
	"well_track/internal/domain/model"
	"well_track/internal/repository"
)

type AnswerUseCase interface {
	SaveAnswer(userID model.UserID, rating int, comment string) error
}

type answerUseCase struct {
	answerRepo repository.AnswerRepository
	log        *zerolog.Logger
}

func NewAnswerUseCase(aRepo repository.AnswerRepository, logger *zerolog.Logger) AnswerUseCase {
	return &answerUseCase{
		answerRepo: aRepo,
		log:        logger,
	}
}

func (uc *answerUseCase) SaveAnswer(userID model.UserID, rating int, comment string) error {
	ans := &model.Answer{
		UserID:    userID,
		Rating:    rating,
		Comment:   comment,
		CreatedAt: time.Now(),
	}
	return uc.answerRepo.Create(ans)
}
