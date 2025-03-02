package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"ddv_loc/pkg/types"
	"ddv_loc/pkg/updater"
)

var checkUpdatesOpts struct {
	input  string
	output string
	format types.FormatEnum
}

var checkUpdatesCmd = &cobra.Command{
	Version: "1.0.0",
	Use:     "check-updates",
	Short:   "Проверка файла/файлов на наличие обновлений",
	Run: func(cmd *cobra.Command, args []string) {
		res, err := updater.CheckUpdates(checkUpdatesOpts.format.String(), checkUpdatesOpts.input, checkUpdatesOpts.output)
		if err != nil {
			fmt.Printf("Ошибка при проверке обновлений: %v\n", err)
			os.Exit(1)
		}

		if res {
			fmt.Printf("Обновления найдены и помещены в %s\n", checkUpdatesOpts.output)
			os.Exit(0)
		}

		fmt.Println("Обновлений нет")
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(checkUpdatesCmd)

	checkUpdatesCmd.Flags().StringVarP(&checkUpdatesOpts.input, "input", "i", "", "Путь до папки с расшифрованным файлом/файлами (required)")
	checkUpdatesCmd.Flags().StringVarP(&checkUpdatesOpts.output, "output", "o", "", "Путь к папке, в которую будет сохранен отчёт (required)")
	checkUpdatesCmd.Flags().VarP(&checkUpdatesOpts.format, "format", "f", "Формат (required)")
	checkUpdatesCmd.MarkFlagRequired("input")
	checkUpdatesCmd.MarkFlagRequired("output")
	checkUpdatesCmd.MarkFlagRequired("format")
}
