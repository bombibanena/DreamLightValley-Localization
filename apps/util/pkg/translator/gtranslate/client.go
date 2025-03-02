package gtranslate

import (
	"github.com/bregydoc/gtranslate"
)

type (
	gtranslateClient struct {
	}
)

func NewClient() *gtranslateClient {
	return &gtranslateClient{}
}

func (c *gtranslateClient) Translate(text string, sourceLang, targetLang string) (string, error) {
	res, err := gtranslate.TranslateWithParams(
		text,
		gtranslate.TranslationParams{
			From: sourceLang,
			To:   targetLang,
		},
	)
	if err != nil {
		return "", err
	}

	return res, nil
}
