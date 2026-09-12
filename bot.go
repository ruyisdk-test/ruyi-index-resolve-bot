package main

import (
	"log/slog"
	"ruyi-index-resolve-bot/bot"
	"ruyi-index-resolve-bot/bot/model"
)

func main() {
	config, err := bot.ConfigLoad()
	if err != nil {
		slog.Error("failed to load config:", "error", err)
		return
	}
	err = model.Hello(config)
	if err != nil {
		slog.Error("failed to greet model:", "error", err)
		return
	}

	resp, err := model.Ask("我们第一次见吗")

	if err != nil {
		return
	}

	slog.Info(resp.Message.Text())

	return
}
