package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/NBVTien/vtdict/internal/config"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type Definition struct {
	Word         string   `json:"word"`
	PartOfSpeech string   `json:"part_of_speech"`
	Definition   string   `json:"definition"`
	Example      string   `json:"example"`
	Synonyms     []string `json:"synonyms"`
}

// providerDefaults maps provider name → (baseURL, defaultModel)
var providerDefaults = map[string][2]string{
	"openai": {"https://api.openai.com/v1", "gpt-4o-mini"},
	"claude": {"https://api.anthropic.com/v1", "claude-haiku-4-5-20251001"},
	"ollama": {"http://localhost:11434/v1", "llama3.2"},
	"groq":   {"https://api.groq.com/openai/v1", "llama-3.1-8b-instant"},
}

// envKeys maps provider → env var name for API key
var envKeys = map[string]string{
	"openai": "OPENAI_API_KEY",
	"claude": "ANTHROPIC_API_KEY",
	"groq":   "GROQ_API_KEY",
	"ollama": "", // no key needed
}

func resolveKey(cfg config.Config) string {
	// config file takes priority
	if cfg.AIKey != "" {
		return cfg.AIKey
	}
	// fall back to env var
	if envVar, ok := envKeys[cfg.AIProvider]; ok && envVar != "" {
		return os.Getenv(envVar)
	}
	return ""
}

func resolveBaseURL(cfg config.Config) string {
	if cfg.AIBaseURL != "" {
		return cfg.AIBaseURL
	}
	if d, ok := providerDefaults[cfg.AIProvider]; ok {
		return d[0]
	}
	return "https://api.openai.com/v1"
}

func resolveModel(cfg config.Config) string {
	if cfg.AIModel != "" {
		return cfg.AIModel
	}
	if d, ok := providerDefaults[cfg.AIProvider]; ok {
		return d[1]
	}
	return "gpt-4o-mini"
}

func LookupWord(word string) (*Definition, error) {
	cfg := config.Get()
	key := resolveKey(cfg)
	baseURL := resolveBaseURL(cfg)
	model := resolveModel(cfg)

	opts := []option.RequestOption{
		option.WithBaseURL(baseURL),
	}
	if key != "" {
		opts = append(opts, option.WithAPIKey(key))
	}

	client := openai.NewClient(opts...)

	prompt := fmt.Sprintf(`Define the English word "%s". Return ONLY valid JSON with this exact shape:
{
  "word": "%s",
  "part_of_speech": "<noun|verb|adjective|adverb|etc>",
  "definition": "<clear definition>",
  "example": "<one natural example sentence>",
  "synonyms": ["<synonym1>", "<synonym2>"]
}
No markdown, no explanation, just the JSON object.`, word, word)

	resp, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Model: model,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("AI lookup failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("AI returned no response")
	}

	var def Definition
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &def); err != nil {
		return nil, fmt.Errorf("AI response parse failed: %w", err)
	}
	return &def, nil
}
