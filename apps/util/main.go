package main

import (
	"ddv_loc/config"
	"ddv_loc/pkg/app"
	"ddv_loc/pkg/cmd"
)

const (
	version = "2.0.0" // Version
)

func main() {
	app.Config = config.New()

	cmd.Execute()
}
