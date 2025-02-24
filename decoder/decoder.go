package decoder

import (
	"fmt"
	"os"
)

func Decode(format string, inPath string, outPath string) error {
	// Read encoded files
	fileData, err := readEncodedFiles(inPath)
	if err != nil {
		fmt.Printf("Ошибка при обработке входной папки: %v\n", err)
		os.Exit(1)
	}

	// Generate output file
	switch format {
	case "json":
		err = generateJSON(fileData, outPath)
	case "csv":
		err = generateCSV(fileData, outPath)
	case "raw":
		err = generateRaw(fileData, outPath)
	default:
		fmt.Printf("Неизвестный формат: %v\n", format)
	}

	if err != nil {
		return fmt.Errorf("Ошибка при генерации выходного файла: %v\n", err)
	}

	return nil
}
