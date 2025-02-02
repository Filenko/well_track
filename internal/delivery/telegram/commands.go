package telegram

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"strconv"
	"strings"
	"well_track/internal/domain/model"
)

func (h *BotHandler) handleStartCommand(b *gotgbot.Bot, ctx *ext.Context) error {
	tgID := ctx.EffectiveUser.Id

	user, err := h.UserUC.GetOrCreateUser(model.TelegramID(tgID))
	if err != nil {
		return err
	}

	text := "Привет! Вызови команду /interval <X>, чтобы задать как часто тебе писать :)"
	_, err = b.SendMessage(int64(user.TelegramID), text, nil)
	return err
}

func (h *BotHandler) handleIntervalCommand(b *gotgbot.Bot, ctx *ext.Context) error {

	sublog := h.log.With().Str("function", "telegram.handleIntervalCommand").Logger()

	chatID := ctx.EffectiveMessage.Chat.Id
	tgUserID := ctx.EffectiveUser.Id
	text := ctx.EffectiveMessage.Text

	parts := strings.Fields(text)

	if len(parts) < 2 {
		_, _ = b.SendMessage(ctx.EffectiveChat.Id, "Пожалуйста, введите число (пока не поздно) после /interval. Пример: /interval 30", nil)
		return nil
	}

	intervalStr := parts[1]
	interval, err := strconv.Atoi(intervalStr)
	if err != nil {
		_, _ = b.SendMessage(ctx.EffectiveChat.Id, "Не удалось понять число (мне каж тут ты виноват ес честно). Убедитесь, что вы ввели корректное целое число.", nil)
		return err
	}

	user, uErr := h.UserUC.GetOrCreateUser(model.TelegramID(tgUserID))
	if uErr != nil {
		sublog.Error().Err(uErr).Msg("Error getting user")
		return uErr
	}
	if err := h.ScheduleUC.SetSchedule(user, interval); err != nil {
		sublog.Error().Err(uErr).Msg("SetSchedule err")
		b.SendMessage(chatID, "Ошибка при сохранении расписания (тут тоже ты виноват мне каж)", nil)
		return err
	}
	b.SendMessage(chatID, "Отлично! Я запомнил, буду напоминать каждые минуты: "+text, nil)

	return nil
}

// handleHelpCommand — общий
func handleHelpCommand(b *gotgbot.Bot, ctx *ext.Context) error {
	chatID := ctx.EffectiveChat.Id
	text := "Помощь:\n" +
		"/start - начать\n" +
		"/help - справка\n" +
		"Или введите число, если хотите задать интервал."
	_, err := b.SendMessage(chatID, text, nil)
	return err
}
