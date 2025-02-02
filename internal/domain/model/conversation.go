package model

import "time"

type ConversationState string

const (
	StateNone           ConversationState = "none"
	StateWaitingRating  ConversationState = "waiting_for_rating"
	StateWaitingComment ConversationState = "waiting_for_comment"
)

type Conversation struct {
	UserID    UserID            `json:"user_id"`
	State     ConversationState `json:"state"`
	Payload   map[string]string `json:"payload"`
	UpdatedAt time.Time         `json:"updated_at"`
}

func (c *Conversation) IsActive() bool {
	return c.State != StateNone
}
