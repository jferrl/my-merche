package bot

import "github.com/yanzay/tbot/v2"

type Handler func(m *tbot.Message)

func WithAuthorizeHandler(c *tbot.Client) (string, Handler) {
	return "/autorize", func(m *tbot.Message) {
		c.SendMessage(m.Chat.ID, "Open https://my-merche.herokuapp.com/login/mercedes")
	}
}
