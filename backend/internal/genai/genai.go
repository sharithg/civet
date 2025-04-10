package genai

import (
	"context"
	"encoding/json"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/sharithg/civet/internal/config"
)

type OpenAi struct {
	client   openai.Client
	cacheDir string
	Config   *config.Config
}

func NewOpenAiClient(config *config.Config) OpenAi {
	client := openai.NewClient(
		option.WithAPIKey(config.OpenAIAPIKey),
	)
	cacheDir := "cache/openai"
	_ = os.MkdirAll(cacheDir, os.ModePerm)

	return OpenAi{
		client:   client,
		cacheDir: cacheDir,
	}
}

func JsonChat[T any](ctx context.Context, o *OpenAi, prompt string, input string, schemaName string, schema interface{}) (T, error) {
	var zero T

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        schemaName,
		Description: openai.String(prompt),
		Schema:      schema,
		Strict:      openai.Bool(true),
	}

	chat, err := o.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(input),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				JSONSchema: schemaParam,
			},
		},
		Model: openai.ChatModelGPT4oMini,
	})
	if err != nil {
		return zero, err
	}

	var result T
	err = json.Unmarshal([]byte(chat.Choices[0].Message.Content), &result)
	if err != nil {
		return zero, err
	}

	return result, nil
}
