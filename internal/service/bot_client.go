package service

import (
	bot "gopkg.in/telegram-bot-api.v4"
)

type BotCmd struct {
	cmd *bot.BotAPI
}

func NewBotCmd(cmd *bot.BotAPI) *BotCmd {
	return &BotCmd{cmd: cmd}
}

func (c BotCmd) Send(msg bot.MessageConfig) (bot.Message, error) {
	return c.cmd.Send(msg)
}

func (c BotCmd) ListenForWebhook(webhook string) bot.UpdatesChannel {
	return c.cmd.ListenForWebhook(webhook)
}