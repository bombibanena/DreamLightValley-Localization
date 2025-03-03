package config

type (
	Config struct {
		API API `mapstructure:"api"`
	}

	API struct {
		URL   string `mapstructure:"url"`
		Token string `mapstructure:"token"`
	}
)
