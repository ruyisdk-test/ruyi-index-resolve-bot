package main

import (
	"log/slog"

	resolvebot "github.com/ruyisdk-test/ruyi-index-resolve-bot/bot"
)

func main() {
	config, err := resolvebot.ConfigLoad()
	if err != nil {
		slog.Error("failed to load config:", "error", err)
		return
	}
	err = resolvebot.ModelHello(config)
	if err != nil {
		slog.Error("failed to greet model:", "error", err)
		return
	}

	resp, err := resolvebot.ModelAsk("我们第一次见吗")

	if err != nil {
		return
	}

	slog.Info(resp.Message.Text())

	return
}
