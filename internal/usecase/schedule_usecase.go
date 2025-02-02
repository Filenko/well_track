package usecase

import (
	"encoding/json"
	"github.com/rs/zerolog"
	"time"
	"well_track/internal/domain/model"
	"well_track/internal/infrastructure/queue"
	"well_track/internal/repository"
)

// ScheduleUseCase отвечает за управление расписаниями и отправку напоминаний
type ScheduleUseCase interface {
	SetSchedule(user *model.User, intervalHours int) error
	SendReminder(userID model.UserID) error
}

type scheduleUseCase struct {
	scheduleRepo repository.ScheduleRepository
	userRepo     repository.UserRepository
	convRepo     repository.ConversationStateRepository
	producer     queue.Producer
	telegramSvc  model.TelegramService
	log          *zerolog.Logger
}

func NewScheduleUseCase(
	sRepo repository.ScheduleRepository,
	uRepo repository.UserRepository,
	cRepo repository.ConversationStateRepository,
	producer queue.Producer,
	telegramSvc model.TelegramService,
	logger *zerolog.Logger,
) ScheduleUseCase {
	return &scheduleUseCase{
		scheduleRepo: sRepo,
		userRepo:     uRepo,
		producer:     producer,
		telegramSvc:  telegramSvc,
		convRepo:     cRepo,
		log:          logger,
	}
}

func (uc *scheduleUseCase) SetSchedule(user *model.User, intervalHours int) error {

	sublog := uc.log.With().Str("function", "usecase.SetSchedule").Logger()

	sched, err := uc.scheduleRepo.GetByUserID(user.ID)

	if err != nil {
		sublog.Error().Err(err).Msg("error getting schedule")
		return err
	}

	schedJson, err := json.Marshal(sched)
	if err != nil {
		sublog.Debug().Err(err).Msg("Error marshalling schedule")
	} else {
		sublog.Debug().RawJSON("sched", schedJson).Msg("Get schedule")
	}

	now := time.Now()
	if sched == nil {
		sched = &model.Schedule{
			UserID:               user.ID,
			IntervalMinutes:      intervalHours,
			LastNotificationTime: now,
		}
	} else {
		sched.IntervalMinutes = intervalHours
		sched.LastNotificationTime = now
	}

	if err := uc.scheduleRepo.Upsert(sched); err != nil {
		sublog.Error().Err(err).Msg("Error updating schedule")
		return err
	}

	delay := time.Duration(intervalHours) * time.Minute
	msg := queue.ReminderMessage{
		UserID: int64(user.ID),
		Action: "SendReminder",
	}
	if err := uc.producer.PublishReminder(msg, delay); err != nil {
		sublog.Error().Err(err).Msg("Failed to publish reminder")
		return err
	}
	reminderMsgJson, err := json.Marshal(sched)
	if err != nil {
		sublog.Debug().Err(err).Msg("Error marshalling reminder")
	} else {
		sublog.Debug().RawJSON("reminder", reminderMsgJson).Dur("delay", delay).Msg("Get reminder")
	}

	return nil
}

func (uc *scheduleUseCase) SendReminder(userID model.UserID) error {

	sublog := uc.log.With().Str("function", "usecase.SendReminder").Int64("user_id", int64(userID)).Logger()

	state, err := uc.convRepo.GetState(userID)
	if err != nil {
		sublog.Debug().Err(err).Msg("Error getting state")
	}

	isAnsweredToPrevReminder := true

	if state == model.StateWaitingRating {
		sublog.Debug().Msg("Need to SendReminder but user didnt ask to previous")
		isAnsweredToPrevReminder = false
	}

	sched, err := uc.scheduleRepo.GetByUserID(userID)

	if err != nil {
		sublog.Debug().Err(err).Msg("error getting schedule")
		return err
	}

	schedJson, err := json.Marshal(sched)
	if err != nil {
		sublog.Debug().Err(err).Msg("Error marshalling schedule")
	} else {
		sublog.Debug().RawJSON("sched", schedJson).Msg("Get schedule")
	}

	if err != nil {
		return err
	}

	if time.Now().Before(sched.NextNotification()) {
		sublog.Debug().Msg("Skip schedule notification")
		return nil
	}

	sched.LastNotificationTime = time.Now()
	if err := uc.scheduleRepo.Upsert(sched); err != nil {
		sublog.Debug().Err(err).Msg("Error updating schedule")
		return err
	}

	if isAnsweredToPrevReminder {
		if err := uc.telegramSvc.SendMessageToUserByUserID(userID, "Как настроение от 1 до 5?"); err != nil {
			sublog.Debug().Err(err).Msg("Error sending reminder")
		}

		err = uc.convRepo.SetState(userID, model.StateWaitingRating)
		if err != nil {
			sublog.Debug().Err(err).Msg("Error setting state")
			return err
		}
	}

	delay := time.Duration(sched.IntervalMinutes) * time.Minute
	msg := queue.ReminderMessage{
		UserID: int64(userID),
		Action: "SendReminder",
	}

	if err := uc.producer.PublishReminder(msg, delay); err != nil {
		sublog.Error().Err(err).Msg("Failed to re-schedule reminder")
	}

	msgJson, err := json.Marshal(msg)
	if err != nil {
		sublog.Debug().Err(err).Msg("Error marshalling reminder")
	} else {
		sublog.Debug().RawJSON("reminder", msgJson).Msg("Re-schedule reminder")
	}

	return nil
}
