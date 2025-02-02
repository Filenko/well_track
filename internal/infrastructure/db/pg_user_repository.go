package db

import (
	"database/sql"
	"github.com/rs/zerolog"
	"well_track/internal/domain/model"
	"well_track/internal/repository"
)

type pgUserRepository struct {
	db  *sql.DB
	log *zerolog.Logger
}

func NewPgUserRepository(db *sql.DB, logger *zerolog.Loggerw) repository.UserRepository {
	return &pgUserRepository{
		db:  db,
		log: logger,
	}
}

// пример таблицы: users(id bigserial primary key, telegram_id bigint, created_at timestamptz)
func (r *pgUserRepository) GetByTelegramID(tgID model.TelegramID) (*model.User, error) {
	row := r.db.QueryRow(`
        SELECT id, telegram_id, created_at FROM users WHERE telegram_id = $1
    `, tgID)
	var u model.User
	var id int64
	if err := row.Scan(&id, &u.TelegramID, &u.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	u.ID = model.UserID(id)
	return &u, nil
}

func (r *pgUserRepository) Create(user *model.User) error {
	row := r.db.QueryRow(`
        INSERT INTO users (telegram_id, created_at)
        VALUES ($1, $2)
        RETURNING id
    `, user.TelegramID, user.CreatedAt)
	var newID int64
	if err := row.Scan(&newID); err != nil {
		return err
	}
	user.ID = model.UserID(newID)
	return nil
}

func (r *pgUserRepository) GetByID(userID model.UserID) (*model.User, error) {
	row := r.db.QueryRow(`
        SELECT id, telegram_id, created_at FROM users WHERE id = $1
    `, userID)
	var u model.User
	var id int64
	if err := row.Scan(&id, &u.TelegramID, &u.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	u.ID = model.UserID(id)
	return &u, nil
}
