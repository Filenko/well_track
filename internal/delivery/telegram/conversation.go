package telegram

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"well_track/internal/domain/model"
)

// handleConversationMessage - когда пользователь вводит текст (не команду)
func (h *BotHandler) handleConversationMessage(b *gotgbot.Bot, ctx *ext.Context) error {

	chatID := ctx.EffectiveChat.Id
	tgUserID := ctx.EffectiveUser.Id
	text := ctx.EffectiveMessage.Text

	sublog := h.log.With().
		Str("func", "handleConversationMessage").
		Int("chatID", int(chatID)).
		Int("tgUserID", int(tgUserID)).
		Logger()

	user, err := h.UserUC.GetOrCreateUser(model.TelegramID(tgUserID))
	if err != nil {
		sublog.Error().Err(err).Msg("GetOrCreateUser error")
		return nil
	}

	reply, err := h.Conversation.ProcessMessage(user.ID, text)
	if err != nil {
		sublog.Error().Err(err).Msg("Conversation error:")
		return nil
	}
	if reply != "" {
		b.SendMessage(chatID, reply, nil)
	}
	return nil
}
