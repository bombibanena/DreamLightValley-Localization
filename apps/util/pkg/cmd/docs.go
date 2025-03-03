package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Generate documentation for CLI commands",
	Run: func(cmd *cobra.Command, args []string) {
		err := doc.GenMarkdownTree(rootCmd, "./docs")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Documentation generated in ./docs")
	},
}

func init() {
	rootCmd.AddCommand(docsCmd)
}
