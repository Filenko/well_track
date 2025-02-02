package usecase

import (
	"github.com/rs/zerolog"
	"time"
	"well_track/internal/domain/model"
	"well_track/internal/repository"
)

type UserUseCase interface {
	GetOrCreateUser(tgID model.TelegramID) (*model.User, error)
	GetUserById(userID model.UserID) (*model.User, error)
}

type userUseCase struct {
	userRepo repository.UserRepository
	log      *zerolog.Logger
}

func NewUserUseCase(uRepo repository.UserRepository, logger *zerolog.Logger) UserUseCase {
	return &userUseCase{
		userRepo: uRepo,
		log:      logger,
	}
}

func (uc *userUseCase) GetOrCreateUser(tgID model.TelegramID) (*model.User, error) {

	sublog := uc.log.With().Str("function", "usecase.GetOrCreateUser").Int64("telegram_id", int64(tgID)).Logger()

	user, err := uc.userRepo.GetByTelegramID(tgID)
	if err != nil {
		sublog.Debug().Err(err).Msg("Failed to get user by telegram id")
		return nil, err
	}
	if user == nil {
		sublog.Debug().Msg("User not found by telegram id, creating new one")
		user = &model.User{
			TelegramID: tgID,
			CreatedAt:  time.Now(),
		}
		if err := uc.userRepo.Create(user); err != nil {
			sublog.Debug().Err(err).Msg("Failed to create user")
			return nil, err
		}
	}
	return user, nil
}

func (uc *userUseCase) GetUserById(userID model.UserID) (*model.User, error) {
	user, err := uc.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
