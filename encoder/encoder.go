package encoder

import (
	"fmt"

	"ddv_loc/models"
)

func Encode(format string, inPath string, outPath string) error {
	var fileData []models.FileData
	var err error

	switch format {
	case "json":
		fileData, err = readFromJSON(inPath)
	case "csv":
		fileData, err = readFromCSV(inPath)
	case "raw":
		fileData, err = readFromRaw(inPath)
	default:
		fmt.Printf("Неизвестный формат: %v\n", format)
	}

	if err != nil {
		return fmt.Errorf("Ошибка при чтении файла %s: %v\n", inPath, err)
	}

	if err := generateEncodedFiles(fileData, outPath); err != nil {
		return fmt.Errorf("Ошибка при генерации файлов: %v\n", err)
	}

	return nil
}
