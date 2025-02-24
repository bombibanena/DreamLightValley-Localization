package encoder

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"ddv_loc/models"
)

func readFromJSON(filePath string) ([]models.FileData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var fileData []models.FileData
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&fileData); err != nil {
		return nil, err
	}

	return fileData, nil
}

func readFromCSV(filePath string) ([]models.FileData, error) {
	// Read file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при чтении файла: %v", err)
	}
	defer file.Close()

	// Create reader
	reader := csv.NewReader(file)

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("Ошибка при чтении заголовка: %v", err)
	}

	// Check format of header
	expectedHeader := []string{"location", "key", "en", "ru"}
	if !reflect.DeepEqual(header, expectedHeader) {
		return nil, fmt.Errorf("Неверный формат заголовка: ожидается %v, получено %v", expectedHeader, header)
	}

	fileDataMap := make(map[string]*models.FileData)

	// Read records
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Ошибка при чтении строки: %v", err)
		}

		// Parse record
		location := record[0]
		key := record[1]
		en := record[2]
		ru := record[3]

		// Check location
		if _, exists := fileDataMap[location]; !exists {
			fileDataMap[location] = &models.FileData{
				Location:   location,
				Dictionary: []models.KeyValue{},
			}
		}

		// fill map
		fileDataMap[location].Dictionary = append(fileDataMap[location].Dictionary, models.KeyValue{
			Key: key,
			Loc: models.KeyValueLoc{
				En: en,
				Ru: ru,
			},
		})
	}

	// Convert map to slice
	var fileDataList []models.FileData
	for _, fd := range fileDataMap {
		fileDataList = append(fileDataList, *fd)
	}

	return fileDataList, nil
}

func readFromRaw(folder string) ([]models.FileData, error) {
	var fileData []models.FileData

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		fileName := info.Name()
		dir := filepath.Dir(path)

		if err != nil {
			return err
		}

		// Skip folders
		if info.IsDir() {
			return nil
		}

		// Skip not .txt files
		if filepath.Ext(info.Name()) != ".txt" {
			return nil
		}

		// Read file
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("--> %s - Error [%s]\n", fileName, dir)
			return fmt.Errorf("Ошибка при чтении файла %s: %v", path, err)
		}

		fmt.Printf("--> %s [%s]\n", fileName, dir)

		// Parse encoded file
		keyValues, err := readTxtFile(data)
		if err != nil {
			fmt.Printf("    Decode - Error\n")
			return fmt.Errorf("Ошибка при разборе файла %s: %v", path, err)
		}

		fmt.Printf("    Decode - OK\n")

		// Fill data
		fileData = append(fileData, models.FileData{
			Location:   strings.Replace(filepath.ToSlash(path), filepath.ToSlash(folder), "", 1),
			Dictionary: keyValues,
		})

		return nil
	})

	return fileData, err
}

func readTxtFile(data []byte) ([]models.KeyValue, error) {
	var keyValues []models.KeyValue

	// Parse file
	lines := strings.Split(string(data), "\n")

	for i := 0; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "TAG: ") {
			key := strings.TrimPrefix(lines[i], "TAG: ")
			if i+1 < len(lines) {
				value := strings.TrimSpace(lines[i+1])

				keyValues = append(keyValues, models.KeyValue{
					Key: key,
					Loc: models.KeyValueLoc{
						En: value,
						Ru: value,
					},
				})

				i++ // Skip next line
			}
		}
	}

	return keyValues, nil
}
