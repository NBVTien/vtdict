package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Lang      string `toml:"lang"`
	Translate bool   `toml:"translate"`
	Phonetic  bool   `toml:"phonetic"`
	Examples  bool   `toml:"examples"`
	POS       bool   `toml:"pos"`

	AIFallback bool   `toml:"ai-fallback"`
	AIProvider string `toml:"ai-provider"` // openai | claude | ollama | groq | <any openai-compat>
	AIKey      string `toml:"ai-key"`
	AIModel    string `toml:"ai-model"`   // optional, uses provider default if empty
	AIBaseURL  string `toml:"ai-base-url"` // optional, for custom endpoints
}

var defaults = Config{
	Lang:       "vi",
	Translate:  true,
	Phonetic:   true,
	Examples:   true,
	POS:        true,
	AIFallback: false,
	AIProvider: "openai",
	AIKey:      "",
	AIModel:    "",
	AIBaseURL:  "",
}

var configPath string
var current Config

func Load() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	path := filepath.Join(home, ".config", "vtdict", "config.toml")
	configPath = path
	current = defaults

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return Save()
	}
	if err != nil {
		return err
	}
	_, err = toml.Decode(string(data), &current)
	return err
}

func Save() error {
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return err
	}

	var buf bytes.Buffer
	buf.WriteString("# vtdict configuration\n")
	buf.WriteString("# Edit this file directly or run: vtdict config\n\n")
	buf.WriteString(fmt.Sprintf("%-12s = %q\n", "lang", current.Lang))
	buf.WriteString(fmt.Sprintf("%-12s = %v\n", "translate", current.Translate))
	buf.WriteString(fmt.Sprintf("%-12s = %v\n", "phonetic", current.Phonetic))
	buf.WriteString(fmt.Sprintf("%-12s = %v\n", "examples", current.Examples))
	buf.WriteString(fmt.Sprintf("%-12s = %v\n\n", "pos", current.POS))
	buf.WriteString("# AI fallback — used when word not found or missing examples\n")
	buf.WriteString(fmt.Sprintf("%-12s = %v\n", "ai-fallback", current.AIFallback))
	buf.WriteString(fmt.Sprintf("%-12s = %q\n", "ai-provider", current.AIProvider))
	buf.WriteString(fmt.Sprintf("%-12s = %q\n", "ai-key", current.AIKey))
	buf.WriteString(fmt.Sprintf("%-12s = %q\n", "ai-model", current.AIModel))
	buf.WriteString(fmt.Sprintf("%-12s = %q\n", "ai-base-url", current.AIBaseURL))

	return os.WriteFile(configPath, buf.Bytes(), 0644)
}

func Get() Config      { return current }
func Set(c Config)     { current = c }
func Defaults() Config { return defaults }
func Path() string     { return configPath }

func Reset() error {
	current = defaults
	return Save()
}
