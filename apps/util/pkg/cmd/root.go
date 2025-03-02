package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Version:           "2.0.0",
	Use:               ".\\ddv_loc.exe",
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true, DisableNoDescFlag: true},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
}
