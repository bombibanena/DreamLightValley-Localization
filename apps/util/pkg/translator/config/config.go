package config

import (
	"encoding/json"
	"os"

	deeplxconfig "ddv_loc/pkg/translator/deeplxapi/config"
	llmconfig "ddv_loc/pkg/translator/llm/config"
)

type (
	Config struct {
		Clients     Clients     `mapstructure:"clients"`
		Current     Client      `mapstructure:"current"`
		Translation Translation `mapstructure:"translation"`
	}

	Clients struct {
		DeeplxAPI deeplxconfig.Config `mapstructure:"deeplx_api"`
		LLM       llmconfig.Config    `mapstructure:"llm"`
	}

	Translation struct {
		TextSize       int    `mapstructure:"text_size"`
		BatchSize      int    `mapstructure:"batch_size"`
		Parts          int    `mapstructure:"parts"`
		Separator      string `mapstructure:"separator"`
		ExceptionsPath string `mapstructure:"exceptions_path"`

		exceptions map[string]string
	}
)

type Client string

const (
	DeeplxAPI     Client = "deeplxAPI"
	DeeplxPackage Client = "deeplxPackage"
	GTranslate    Client = "gtranslate"
	LLM           Client = "llm"
)

func (t *Translation) GetExceptions() map[string]string {
	if t.exceptions == nil {
		file, err := os.Open(t.ExceptionsPath)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		var exceptions map[string]string
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&exceptions); err != nil {
			if err != nil {
				panic(err)
			}
		}

		t.exceptions = exceptions
	}

	return t.exceptions
}
