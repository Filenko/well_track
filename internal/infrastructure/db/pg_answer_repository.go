package db

import (
	"database/sql"
	"github.com/rs/zerolog"
	"well_track/internal/domain/model"
	"well_track/internal/repository"
)

type pgAnswerRepository struct {
	db  *sql.DB
	log *zerolog.Logger
}

func NewPgAnswerRepository(db *sql.DB, logger *zerolog.Logger) repository.AnswerRepository {
	return &pgAnswerRepository{
		db:  db,
		log: logger,
	}
}

func (r *pgAnswerRepository) Create(a *model.Answer) error {
	row := r.db.QueryRow(`
        INSERT INTO answers(user_id, rating, comment, created_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `, a.UserID, a.Rating, a.Comment, a.CreatedAt)
	var newID int64
	if err := row.Scan(&newID); err != nil {
		return err
	}
	a.ID = model.AnswerID(newID)
	return nil
}
