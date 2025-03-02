package progress

import (
	"fmt"

	"github.com/ermanimer/progress_bar"
)

type (
	ProgressConfig struct {
		Message string
	}
)

func CreateProgressBar(total float64, cfg ProgressConfig) *progress_bar.ProgressBar {
	fmt.Printf("%s:\n", cfg.Message)
	fmt.Print("  ")

	return progress_bar.DefaultProgressBar(total)
}
