package model

import (
	"context"
	"fmt"
	"log/slog"
	"ruyi-index-resolve-bot/bot"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/compat_oai"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

var botGenKit *genkit.Genkit = nil
var botModel *ai.ModelRef = nil

func Hello(config *bot.Config) error {

	ctx := context.Background()

	g := genkit.Init(ctx, genkit.WithPlugins(&compat_oai.OpenAICompatible{
		Provider: config.Provider,
		APIKey:   config.ApiKey,
		BaseURL:  config.BaseUrl,
		Opts: []option.RequestOption{
			option.WithHeader("Custom-Header", "value"),
		},
	}))

	modelName := fmt.Sprintf("%s/%s", config.Provider, config.ModelName)
	model := ai.NewModelRef(
		modelName,
		&openai.ChatCompletionNewParams{
			Temperature:         openai.Float(0.7),
			MaxCompletionTokens: openai.Int(1024),
		})

	resp, err := genkit.Generate(ctx, g,
		ai.WithModel(model),
		ai.WithPrompt("用一行文字简单介绍你自己"),
	)

	if err != nil {
		return err
	}

	slog.Info(resp.Message.Text())

	botGenKit = g
	botModel = &model

	return nil
}

func Ask(msg string) (*ai.ModelResponse, error) {
	return genkit.Generate(
		context.Background(),
		botGenKit,
		ai.WithModel(botModel),
		ai.WithPrompt(msg),
	)
}

type PackageResults struct {
	Results []PackageResult `json:"results"`
}

type PackageResult struct {
	Name    string   `json:"name"`
	Version string   `json:"version"`
	URLs    []string `json:"urls"`
}

func AskData(msg string) (*PackageResults, *ai.ModelResponse, error) {
	return genkit.GenerateData[PackageResults](
		context.Background(),
		botGenKit,
		ai.WithModel(botModel),
		ai.WithPrompt(msg),
	)
}
