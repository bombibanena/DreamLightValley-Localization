package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"ddv_loc/pkg/encoder"
	"ddv_loc/pkg/types"
)

var encodeOpts struct {
	input  string
	output string
	format types.FormatEnum
}

var encodeCmd = &cobra.Command{
	Version: "1.0.0",
	Use:     "encode",
	Short:   "Шифрование в .locbin файлы",
	Run: func(cmd *cobra.Command, args []string) {
		if err := encoder.Encode(encodeOpts.format.String(), encodeOpts.input, encodeOpts.output); err != nil {
			fmt.Printf("Ошибка при шифровании: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Шифрование завершено\n")
	},
}

func init() {
	rootCmd.AddCommand(encodeCmd)

	encodeCmd.Flags().StringVarP(&encodeOpts.input, "input", "i", "", "Путь до папки с расшифрованным файлом/файлами (required)")
	encodeCmd.Flags().StringVarP(&encodeOpts.output, "output", "o", "", "Путь к папке, в которую будут сохранены зашифрованные файлы (required)")
	encodeCmd.Flags().VarP(&encodeOpts.format, "format", "f", "Формат (required)")
	encodeCmd.MarkFlagRequired("input")
	encodeCmd.MarkFlagRequired("output")
	encodeCmd.MarkFlagRequired("format")
}
