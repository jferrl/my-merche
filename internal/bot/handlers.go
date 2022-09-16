package bot

import (
	"fmt"

	"github.com/yanzay/tbot/v2"
)

type Handler func(m *tbot.Message)

func WithStartHandler(c *tbot.Client) (string, Handler) {
	return "/start", func(m *tbot.Message) {
		c.SendChatAction(m.Chat.ID, tbot.ActionTyping)

		cmds, err := c.GetMyCommands()
		if err != nil {
			c.SendMessage(m.Chat.ID, "Error getting bot commands!")
			return
		}

		var res string
		for _, cmd := range *cmds {
			res += fmt.Sprintf("%s: %s\n", cmd.Command, cmd.Description)
		}

		c.SendMessage(m.Chat.ID, res)
	}
}

func WithLoginHandler(c *tbot.Client) (string, Handler) {
	return "/login", func(m *tbot.Message) {
		c.SendMessage(m.Chat.ID, "Open https://my-merche.herokuapp.com/login/mercedes")
	}
}
