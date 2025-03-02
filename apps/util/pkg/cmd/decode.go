package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"ddv_loc/pkg/decoder"
	"ddv_loc/pkg/types"
)

var decodeOpts struct {
	input  string
	output string
	format types.FormatEnum
}

var decodeCmd = &cobra.Command{
	Version: "1.0.0",
	Use:     "decode",
	Short:   "Расшифровка .locbin файлов",
	Run: func(cmd *cobra.Command, args []string) {
		if err := decoder.Decode(decodeOpts.format.String(), decodeOpts.input, decodeOpts.output); err != nil {
			fmt.Printf("Ошибка при расшифровке: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Расшифровка завершена")
	},
}

func init() {
	rootCmd.AddCommand(decodeCmd)

	decodeCmd.Flags().StringVarP(&decodeOpts.input, "input", "i", "", "Путь до папки с .locbin файлами (required)")
	decodeCmd.Flags().StringVarP(&decodeOpts.output, "output", "o", "", "Путь к папке, в которую будут сохранены расшифрованные файлы (required)")
	decodeCmd.Flags().VarP(&decodeOpts.format, "format", "f", "Формат (required)")
	decodeCmd.MarkFlagRequired("input")
	decodeCmd.MarkFlagRequired("output")
	decodeCmd.MarkFlagRequired("format")
}
