package deeplxpack

import (
	"errors"

	"github.com/xiaoxuan6/deeplx"
)

type (
	deeplxPackageClient struct {
	}
)

func NewClient() *deeplxPackageClient {
	return &deeplxPackageClient{}
}

func (c *deeplxPackageClient) Translate(text string, sourceLang, targetLang string) (string, error) {
	res := deeplx.Translate(text, sourceLang, targetLang)
	if res.Code != 200 || res.Msg != "success" {
		return "", errors.New(res.Msg)
	}

	return res.Data, nil
}
