package telegram

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	"github.com/rs/zerolog"
	"strconv"

	"well_track/internal/domain/model"
	"well_track/internal/usecase"
)

// TelegramService - интерфейс, который UseCase может вызывать

type BotHandler struct {
	Bot          *gotgbot.Bot
	Updater      *ext.Updater
	Conversation usecase.ConversationUseCase
	UserUC       usecase.UserUseCase
	ScheduleUC   usecase.ScheduleUseCase
	log          *zerolog.Logger
}

// Гарантируем, что BotHandler реализует TelegramService
var _ model.TelegramService = (*BotHandler)(nil)

func (h *BotHandler) SendMessageToUserByTelegramID(userID model.TelegramID, text string) error {

	h.log.With().Str("function", "telegram.SendMessageToUserByTelegramID")

	h.log.Debug().
		Str("telegramID", strconv.FormatInt(int64(userID), 10)).
		Str("text", text).
		Msg("SendMessageToUserByTelegramID called")

	_, sendErr := h.Bot.SendMessage(int64(userID), text, nil)
	return sendErr
}

func (h *BotHandler) SendMessageToUserByUserID(userID model.UserID, text string) error {
	h.log.With().Str("function", "telegram.SendMessageToUserByUserID")

	h.log.Debug().
		Str("telegramID", strconv.FormatInt(int64(userID), 10)).
		Str("text", text).
		Msg("SendMessageToUserByUserID called")

	user, err := h.UserUC.GetUserById(userID)
	if err != nil {
		return err
	}
	chatID := user.TelegramID
	_, sendErr := h.Bot.SendMessage(int64(chatID), text, nil)
	return sendErr
}

func NewBotHandler(
	botToken string,
	convUC usecase.ConversationUseCase,
	userUC usecase.UserUseCase,
	schedUC usecase.ScheduleUseCase,
	logger *zerolog.Logger,
) (*BotHandler, error) {
	logger.With().Str("function", "telegram.NewBotHandler")
	bot, err := gotgbot.NewBot(botToken, nil)

	if err != nil {
		return nil, err
	}

	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			logger.Err(err).Msg("An error occurred while handling update")
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})

	updater := ext.NewUpdater(dispatcher, nil)

	handler := &BotHandler{
		Bot:          bot,
		Updater:      updater,
		Conversation: convUC,
		UserUC:       userUC,
		ScheduleUC:   schedUC,
		log:          logger,
	}

	dispatcher.AddHandler(handlers.NewCommand("start", handler.handleStartCommand))
	dispatcher.AddHandler(handlers.NewCommand("help", handleHelpCommand))
	dispatcher.AddHandler(handlers.NewCommand("interval", handler.handleIntervalCommand))

	dispatcher.AddHandler(handlers.NewMessage(message.Text, handler.handleConversationMessage))

	return handler, nil
}

func (h *BotHandler) StartPolling() error {
	err := h.Updater.StartPolling(h.Bot, &ext.PollingOpts{
		DropPendingUpdates: true,
	})
	if err != nil {
		return err
	}
	h.log.Info().Str("username", h.Bot.User.Username).Msg("Bot @%s started...\n")
	return nil
}
