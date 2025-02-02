package repository

import (
	"well_track/internal/domain/model"
)

type ConversationStateRepository interface {
	GetState(userID model.UserID) (model.ConversationState, error)
	SetState(userID model.UserID, state model.ConversationState) error
	SetPayload(userID model.UserID, data map[string]string) error
	GetPayload(userID model.UserID) (map[string]string, error)
}
