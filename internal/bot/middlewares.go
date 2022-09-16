package bot

import (
	"time"

	"github.com/yanzay/tbot/v2"
)

func WithSecure(allowlist Allowlist, c *tbot.Client) tbot.Middleware {
	return func(h tbot.UpdateHandler) tbot.UpdateHandler {
		return func(u *tbot.Update) {
			if exists, allowed := allowlist[u.Message.From.Username]; !exists || !allowed || u.Message.From.IsBot {
				c.SendChatAction(u.Message.Chat.ID, tbot.ActionTyping)
				time.Sleep(1 * time.Second)
				c.SendMessage(u.Message.Chat.ID, "Not allowed to use the bot!")
				return
			}

			h(u)
		}
	}
}
