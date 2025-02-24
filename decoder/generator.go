package decoder

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ddv_loc/models"
)

func generateJSON(fileData []models.FileData, outPath string) error {
	// Create out directory if it does not exist
	if err := os.MkdirAll(outPath, os.ModePerm); err != nil {
		return err
	}

	// Create out file
	filePath := filepath.Join(outPath, "loc.json")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Generate JSON
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(fileData)
}

func generateCSV(fileData []models.FileData, outPath string) error {
	// Create out directory if it does not exist
	if err := os.MkdirAll(outPath, os.ModePerm); err != nil {
		return err
	}

	// Create out file
	filePath := filepath.Join(outPath, "loc.csv")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Header
	if err := writer.Write([]string{"location", "key", "en", "ru"}); err != nil {
		return err
	}

	// Body
	for _, fd := range fileData {
		for _, kv := range fd.Dictionary {
			record := []string{
				fd.Location,
				kv.Key,
				kv.Loc.En,
				kv.Loc.Ru,
			}
			if err := writer.Write(record); err != nil {
				return err
			}
		}
	}

	return nil
}

func generateRaw(fileData []models.FileData, outPath string) error {
	// Create out directory if it does not exist
	if err := os.MkdirAll(outPath, os.ModePerm); err != nil {
		return err
	}

	for _, fd := range fileData {
		if err := generateTxtFile(fd, outPath); err != nil {
			return fmt.Errorf("Ошибка при записи файла %s: %v", fd, err)
		}
	}

	return nil
}

func generateTxtFile(fileData models.FileData, outFolder string) error {
	filePath := filepath.Join(outFolder, fileData.Location)
	filePath = strings.Replace(filePath, ".locbin", ".txt", 1)

	fmt.Printf("--> %s\n", filePath)

	err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		return fmt.Errorf("Ошибка при создании папки %s: %v", filepath.Dir(filePath), err)
	}

	outFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("Ошибка при создании файла %s: %v", filePath, err)
	}
	defer outFile.Close()

	for _, kv := range fileData.Dictionary {
		outFile.WriteString("TAG: " + kv.Key + "\n")
		outFile.WriteString(kv.Loc.Ru + "\n\n")
	}

	return nil
}
