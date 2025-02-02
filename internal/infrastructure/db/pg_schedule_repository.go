package db

import (
	"database/sql"
	"github.com/rs/zerolog"
	"time"
	"well_track/internal/domain/model"
	"well_track/internal/repository"
)

type pgScheduleRepository struct {
	db  *sql.DB
	log *zerolog.Logger
}

func NewPgScheduleRepository(db *sql.DB, logger *zerolog.Logger) repository.ScheduleRepository {
	return &pgScheduleRepository{
		db:  db,
		log: logger,
	}
}

// предполагаем таблицу schedules(id bigserial, user_id bigint, interval_hours int, last_notification_time timestamptz)
func (r *pgScheduleRepository) GetByUserID(userID model.UserID) (*model.Schedule, error) {
	row := r.db.QueryRow(`
        SELECT id, user_id, interval_hours, last_notification_time
        FROM schedules
        WHERE user_id = $1
    `, userID)
	var s model.Schedule
	var id int64
	var uid int64
	var lnt time.Time
	if err := row.Scan(&id, &uid, &s.IntervalMinutes, &lnt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	s.ID = model.ScheduleID(id)
	s.UserID = model.UserID(uid)
	s.LastNotificationTime = lnt
	return &s, nil
}

func (r *pgScheduleRepository) Upsert(s *model.Schedule) error {
	// upsert по user_id
	// PostgreSQL 9.5+ можно делать ON CONFLICT (user_id) DO UPDATE
	// Для упрощения сделаем SELECT -> INSERT / UPDATE
	existing, err := r.GetByUserID(s.UserID)
	if err != nil {
		return err
	}
	if existing == nil {
		// INSERT
		row := r.db.QueryRow(`
            INSERT INTO schedules(user_id, interval_hours, last_notification_time)
            VALUES ($1, $2, $3)
            RETURNING id
        `, s.UserID, s.IntervalMinutes, s.LastNotificationTime)
		var newID int64
		if err := row.Scan(&newID); err != nil {
			return err
		}
		s.ID = model.ScheduleID(newID)
	} else {
		// UPDATE
		_, err := r.db.Exec(`
            UPDATE schedules
            SET interval_hours = $1, last_notification_time = $2
            WHERE id = $3
        `, s.IntervalMinutes, s.LastNotificationTime, existing.ID)
		if err != nil {
			return err
		}
		s.ID = existing.ID
	}
	return nil
}
