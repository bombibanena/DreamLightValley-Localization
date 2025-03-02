package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"ddv_loc/pkg/translator"
	"ddv_loc/pkg/types"
)

var translateOpts struct {
	input    string
	output   string
	format   types.FormatEnum
	language string
}

var translateCmd = &cobra.Command{
	Version: "1.0.0",
	Use:     "translate",
	Short:   "Перевод расшифрованных файлов",
	Run: func(cmd *cobra.Command, args []string) {
		if err := translator.Translate(
			translateOpts.format.String(),
			translateOpts.input,
			translateOpts.output,
			translateOpts.language,
		); err != nil {
			fmt.Printf("Ошибка при переводе: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Перевод завершен")
	},
}

func init() {
	rootCmd.AddCommand(translateCmd)

	translateCmd.Flags().StringVarP(&translateOpts.input, "input", "i", "", "Путь до файла/файлов, которые нужно перевести (required)")
	translateCmd.Flags().StringVarP(&translateOpts.output, "output", "o", "", "Путь к папке, в которую будут сохранены переведенные файлы (required)")
	translateCmd.Flags().VarP(&translateOpts.format, "format", "f", "Формат (required)")
	translateCmd.Flags().StringVarP(&translateOpts.language, "language", "l", "RU", "Язык")
	translateCmd.MarkFlagRequired("input")
	translateCmd.MarkFlagRequired("output")
	translateCmd.MarkFlagRequired("format")
}
