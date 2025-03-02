package main

import (
	"ddv_loc/config"
	"ddv_loc/pkg/app"
	"ddv_loc/pkg/cmd"
)

func main() {
	app.Config = config.New()

	cmd.Execute()
}
