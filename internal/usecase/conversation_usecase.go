package usecase

import (
	"github.com/rs/zerolog"
	"strconv"
	"strings"
	"well_track/internal/domain/model"
	"well_track/internal/repository"
)

type ConversationUseCase interface {
	ProcessMessage(userID model.UserID, text string) (string, error)
}

type conversationUseCase struct {
	fsmRepo  repository.ConversationStateRepository
	answerUC AnswerUseCase
	log      *zerolog.Logger
}

func NewConversationUseCase(
	fsmRepo repository.ConversationStateRepository,
	answerUC AnswerUseCase,
	logger *zerolog.Logger,
) ConversationUseCase {
	return &conversationUseCase{
		fsmRepo:  fsmRepo,
		answerUC: answerUC,
		log:      logger,
	}
}

func (uc *conversationUseCase) ProcessMessage(userID model.UserID, text string) (string, error) {

	sublog := uc.log.With().Str("function", "usecase.ProcessMessage").Int("userID", int(userID)).Logger()

	state, err := uc.fsmRepo.GetState(userID)
	if err != nil {
		sublog.Error().Err(err).Msg("failed to fetch state")
		return "", err
	}
	sublog.Debug().Str("User state", string(state)).Msg("Get user state")

	switch state {

	case model.StateNone:
		return "Я пока не умею делать ничего, кроме оценки состояния :(", nil

	case model.StateWaitingRating:
		ratingStr := strings.TrimSpace(text)
		rating, parseErr := strconv.Atoi(ratingStr)
		if parseErr != nil || rating < 1 || rating > 5 {
			return "Введите число от 1 до 5.", nil
		}
		ratingPayload := map[string]string{
			"rating": ratingStr,
		}
		err := uc.fsmRepo.SetPayload(userID, ratingPayload)
		if err != nil {
			sublog.Error().Err(err).Msg("Failed to set payload")
			return "", err
		}

		err = uc.fsmRepo.SetState(userID, model.StateWaitingComment)
		if err != nil {
			sublog.Error().Err(err).Msg("failed to set state")
			return "", err
		}
		return "Расскажите, что у вас происходило в последнее время:", nil

	case model.StateWaitingComment:
		comment := text
		ratingPayload, err := uc.fsmRepo.GetPayload(userID)
		if err != nil {
			sublog.Error().Err(err).Msg("failed to get payload")
			return "", err
		}
		rating, err := strconv.Atoi(ratingPayload["rating"])
		if err != nil {
			return "", err
		}

		err = uc.answerUC.SaveAnswer(userID, rating, comment)
		if err != nil {
			sublog.Error().Err(err).Msg("failed to save answer")
			return "", err
		}

		_ = uc.fsmRepo.SetState(userID, model.StateNone)
		return "Спасибо! Записал ваш комментарий.", nil

	default:
		// если что-то неизвестно
		_ = uc.fsmRepo.SetState(userID, model.StateNone)
		return "Сброс. Попробуйте заново.", nil
	}
}
