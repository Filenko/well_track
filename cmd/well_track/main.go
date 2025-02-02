package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"os"
	"well_track/internal/config"
	"well_track/internal/delivery/rabbitmq"
	"well_track/internal/delivery/telegram"
	"well_track/internal/infrastructure/cache"
	"well_track/internal/infrastructure/db"
	"well_track/internal/infrastructure/logger"
	"well_track/internal/infrastructure/queue"
	"well_track/internal/usecase"
)

func main() {

	log := logger.New()

	cfg := config.MustLoad(log)

	if cfg.Telegram.TelegramApiToken == "" {
		log.Fatal().Msg("TELEGRAM_BOT_TOKEN environment variable not set")
	}

	// 2. Инициализация Postgres
	dbConn, err := db.NewPostgresDB(&cfg.Postgres)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to postgres")
	}
	log.Info().Msg("Connected to Postgres")

	redisClient := redis.NewClient(&redis.Options{Addr: cfg.Redis.Address})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatal().Err(err).Msg("Failed to connect Redis")
	} else {
		log.Info().Str("address", cfg.Redis.Address).Msg("Connected to Redis")
	}

	userRepo := db.NewPgUserRepository(dbConn, log)
	scheduleRepo := db.NewPgScheduleRepository(dbConn, log)
	answerRepo := db.NewPgAnswerRepository(dbConn, log)
	conversationStateRepo := cache.NewConversationStateRepositoryRedis(redisClient, log)

	rabbitConn, err := queue.NewRabbitMQConnection(&cfg.RabbitMQ, log)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect Rabbit")
	}

	ch, err := rabbitConn.Channel()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open channel")
	}

	producer, err := queue.NewRabbitProducer(ch, "reminder_exchange", "reminder", log)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to init producer")
	}

	consumer := queue.NewRabbitConsumer(ch, log)

	userUC := usecase.NewUserUseCase(userRepo, log)
	answerUC := usecase.NewAnswerUseCase(answerRepo, log)
	conversationUC := usecase.NewConversationUseCase(conversationStateRepo, answerUC, log)

	botHandler, err := telegram.NewBotHandler(cfg.Telegram.TelegramApiToken, conversationUC, userUC, nil, log)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create bot handler")
	}

	scheduleUC := usecase.NewScheduleUseCase(scheduleRepo, userRepo, conversationStateRepo, producer, botHandler, log)

	botHandler.ScheduleUC = scheduleUC

	schedConsumer := rabbitmq.NewScheduleConsumer(scheduleUC, log)

	if err := consumer.StartConsume("reminder_queue", schedConsumer.HandleReminder); err != nil {
		log.Fatal().Err(err).Msg("Failed to start consumer")
	}

	if err := botHandler.StartPolling(); err != nil {
		log.Fatal().Err(err).Msg("Bot polling failed")
	}

	log.Info().Msg("Bot is running...")

	select {}
}

func envOrDefault(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}
