package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"well_track/internal/domain/model"
	"well_track/internal/repository"

	"github.com/redis/go-redis/v9"
)

type ConversationStateRepositoryRedis struct {
	client *redis.Client
	log    *zerolog.Logger
}

func NewConversationStateRepositoryRedis(client *redis.Client, logger *zerolog.Logger) repository.ConversationStateRepository {
	return &ConversationStateRepositoryRedis{
		client: client,
		log:    logger,
	}
}

func (r *ConversationStateRepositoryRedis) GetState(userID model.UserID) (model.ConversationState, error) {
	ctx := context.Background()

	key := fmt.Sprintf("convstate:%d:state", userID)
	val, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return model.StateNone, nil
	}
	if err != nil {
		return "", err
	}
	return model.ConversationState(val), nil
}

func (r *ConversationStateRepositoryRedis) SetState(userID model.UserID, state model.ConversationState) error {
	ctx := context.Background()

	key := fmt.Sprintf("convstate:%d:state", userID)
	return r.client.Set(ctx, key, string(state), 0).Err()
}

func (r *ConversationStateRepositoryRedis) SetPayload(userID model.UserID, data map[string]string) error {
	ctx := context.Background()

	key := fmt.Sprintf("convstate:%d:payload", userID)

	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, bytes, 0).Err()
}

func (r *ConversationStateRepositoryRedis) GetPayload(userID model.UserID) (map[string]string, error) {
	ctx := context.Background()

	key := fmt.Sprintf("convstate:%d:payload", userID)
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return map[string]string{}, nil
	}
	if err != nil {
		return nil, err
	}

	data := make(map[string]string)
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, err
	}
	return data, nil
}
