package repository

import (
	"well_track/internal/domain/model"
)

type ScheduleRepository interface {
	GetByUserID(userID model.UserID) (*model.Schedule, error)
	Upsert(schedule *model.Schedule) error
}
