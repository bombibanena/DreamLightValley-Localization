package llm

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"ddv_loc/pkg/translator/llm/config"
	pb "ddv_loc/pkg/translator/llm/translation.v1"
)

type (
	llmClient struct {
		cfg config.Config
	}
)

func NewClient(cfg config.Config) *llmClient {
	return &llmClient{
		cfg: cfg,
	}
}

func (c *llmClient) Translate(text string, sourceLang, targetLang string) (string, error) {
	conn, err := grpc.NewClient(c.cfg.Grpc.Target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", err
	}
	defer conn.Close()

	client := pb.NewTranslationServiceClient(conn)
	res, err := client.Translate(
		context.Background(),
		&pb.TranslationRequest{
			Text: text,
		},
	)
	if err != nil {
		return "", err
	}

	return res.Data, nil
}
