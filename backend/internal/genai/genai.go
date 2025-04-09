package genai

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

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

	cacheKey := computeCacheKey(prompt, input, schemaName, schema)
	cachePath := filepath.Join(o.cacheDir, fmt.Sprintf("%s.json", cacheKey))

	if content, err := os.ReadFile(cachePath); err == nil {
		var cached T
		if err := json.Unmarshal(content, &cached); err == nil {
			fmt.Printf("Loaded cached OpenAI response: %s\n", cacheKey)
			return cached, nil
		}
	}

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

	if encoded, err := json.MarshalIndent(result, "", "  "); err == nil {
		_ = os.WriteFile(cachePath, encoded, 0644)
	}

	return result, nil
}

func computeCacheKey(prompt, input, schemaName string, schema interface{}) string {
	data, _ := json.Marshal(struct {
		Prompt     string      `json:"prompt"`
		Input      string      `json:"input"`
		SchemaName string      `json:"schema_name"`
		Schema     interface{} `json:"schema"`
	}{
		Prompt:     prompt,
		Input:      input,
		SchemaName: schemaName,
		Schema:     schema,
	})
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash[:])
}
