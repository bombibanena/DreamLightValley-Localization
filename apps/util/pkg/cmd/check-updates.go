package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"ddv_loc/pkg/types"
	"ddv_loc/pkg/updater"
)

var checkUpdatesOpts struct {
	inputOld string
	inputNew string
	report   string
	format   types.FormatEnum
}

var checkUpdatesCmd = &cobra.Command{
	Version: "1.0.0",
	Use:     "check-updates",
	Short:   "Проверка файла/файлов на наличие обновлений",
	Run: func(cmd *cobra.Command, args []string) {
		res, err := updater.CheckUpdates(checkUpdatesOpts.format.String(), checkUpdatesOpts.inputOld, checkUpdatesOpts.inputNew, checkUpdatesOpts.report)
		if err != nil {
			fmt.Printf("Ошибка при проверке обновлений: %v\n", err)
			os.Exit(1)
		}

		if res {
			fmt.Printf("Обновления найдены и помещены в %s\n", checkUpdatesOpts.report)
			os.Exit(0)
		}

		fmt.Println("Обновлений нет")
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(checkUpdatesCmd)

	checkUpdatesCmd.Flags().StringVarP(&checkUpdatesOpts.inputOld, "input-old", "o", "", "Путь до папки с расшифрованным файлом/файлами (СТАРЫЙ) (required)")
	checkUpdatesCmd.Flags().StringVarP(&checkUpdatesOpts.inputNew, "input-new", "n", "", "Путь до папки с расшифрованным файлом/файлами (НОВЫЙ) (required)")
	checkUpdatesCmd.Flags().StringVarP(&checkUpdatesOpts.report, "report", "r", "", "Путь к папке, в которую будет сохранен отчёт (required)")
	checkUpdatesCmd.Flags().VarP(&checkUpdatesOpts.format, "format", "f", "Формат (required)")
	checkUpdatesCmd.MarkFlagRequired("input-old")
	checkUpdatesCmd.MarkFlagRequired("input-new")
	checkUpdatesCmd.MarkFlagRequired("report")
	checkUpdatesCmd.MarkFlagRequired("format")
}
