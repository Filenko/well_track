package rabbitmq

import (
	"encoding/json"
	"github.com/rs/zerolog"
	"well_track/internal/domain/model"
	"well_track/internal/infrastructure/queue"
	"well_track/internal/usecase"
)

type ScheduleConsumer struct {
	scheduleUC usecase.ScheduleUseCase
	log        *zerolog.Logger
}

func NewScheduleConsumer(suc usecase.ScheduleUseCase, logger *zerolog.Logger) *ScheduleConsumer {
	return &ScheduleConsumer{
		scheduleUC: suc,
		log:        logger,
	}
}

func (c *ScheduleConsumer) HandleReminder(msg queue.ReminderMessage) error {
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		c.log.Debug().Msg("Get new reminder, but error marshalling it.")
	} else {
		c.log.Debug().RawJSON("msg", jsonMsg).Msg("Get new reminder.")
	}

	if msg.Action == "SendReminder" {
		userID := msg.UserID
		if err := c.scheduleUC.SendReminder(model.UserID(userID)); err != nil {
			c.log.Error().Err(err).Msg("Error sending reminder")
			return err
		}
	}
	return nil
}
