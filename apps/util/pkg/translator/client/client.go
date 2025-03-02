package client

import (
	"ddv_loc/pkg/translator/config"
	"ddv_loc/pkg/translator/deeplxapi"
	"ddv_loc/pkg/translator/deeplxpack"
	"ddv_loc/pkg/translator/gtranslate"
	"ddv_loc/pkg/translator/llm"
)

type IClient interface {
	Translate(text string, sourceLang string, targetLang string) (string, error)
}

func GetClient(cfg config.Config) IClient {
	switch cfg.Current {
	case config.DeeplxAPI:
		return deeplxapi.NewClient(cfg.Clients.DeeplxAPI)
	case config.DeeplxPackage:
		return deeplxpack.NewClient()
	case config.GTranslate:
		return gtranslate.NewClient()
	case config.LLM:
		return llm.NewClient(cfg.Clients.LLM)
	default:
		return nil
	}
}
