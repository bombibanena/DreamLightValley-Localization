package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"ddv_loc/pkg/types"
	"ddv_loc/pkg/updater"
)

var patchOpts struct {
	input     string
	output    string
	updates   string
	translate bool
	format    types.FormatEnum
}

var patchCmd = &cobra.Command{
	Version: "1.0.0",
	Use:     "patch",
	Short:   "Обновление файлов",
	Run: func(cmd *cobra.Command, args []string) {
		if err := updater.Patch(patchOpts.format.String(), patchOpts.input, patchOpts.updates, patchOpts.output, patchOpts.translate); err != nil {
			fmt.Printf("Ошибка при обновлении: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Обновление завершено\n")
	},
}

func init() {
	rootCmd.AddCommand(patchCmd)

	patchCmd.Flags().StringVarP(&patchOpts.input, "input", "i", "", "Путь до папки с расшифрованным файлом JSON (required)")
	patchCmd.Flags().StringVarP(&patchOpts.output, "output", "o", "", "Путь к папке, в которую будут сохранены результаты (required)")
	patchCmd.Flags().StringVarP(&patchOpts.updates, "updates", "u", "", "Путь к файлу с обновлениями (required)")
	patchCmd.Flags().VarP(&patchOpts.format, "format", "f", "Формат (required)")
	patchCmd.Flags().BoolVarP(&patchOpts.translate, "translate", "t", false, "Делать перевод")
	patchCmd.MarkFlagRequired("input")
	patchCmd.MarkFlagRequired("output")
	patchCmd.MarkFlagRequired("updates")
	patchCmd.MarkFlagRequired("format")
}
