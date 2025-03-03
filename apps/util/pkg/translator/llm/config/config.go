package config

type (
	Config struct {
		Grpc Grpc `mapstructure:"grpc"`
	}

	Grpc struct {
		Target string `mapstructure:"target"`
	}
)
