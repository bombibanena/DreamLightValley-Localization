package generator

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ddv_loc/pkg/models"
	"ddv_loc/pkg/types"
	"ddv_loc/pkg/utils/progress"
)

func GenerateDecodedJSON(locFile types.LocFile, outPath string) error {
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
	return encoder.Encode(locFile)
}

func GenerateDecodedCSV(locFile types.LocFile, outPath string) error {
	// Create out directory if it does not exist
	if err := os.MkdirAll(outPath, os.ModePerm); err != nil {
		return err
	}

	fileName := "loc.csv"

	// Create out file
	filePath := filepath.Join(outPath, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	pb := progress.CreateProgressBar(float64(len(locFile)), progress.ProgressConfig{Message: fmt.Sprintf("Генерация %s", fileName)})
	pb.Start()

	// Header
	if err := writer.Write([]string{"location", "key", "en", "ru"}); err != nil {
		return err
	}

	// Body
	for _, fd := range locFile {
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

		pb.Update(pb.CurrentValue + 1)
	}

	return nil
}

func GenerateDecodedTxt(locFile types.LocFile, folder string) error {
	// Create out directory if it does not exist
	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		return err
	}

	pb := progress.CreateProgressBar(float64(len(locFile)), progress.ProgressConfig{Message: "Генерация .txt файлов"})
	pb.Start()

	for _, fd := range locFile {
		if err := generateTxtFile(fd, folder); err != nil {
			return err
		}

		pb.Update(pb.CurrentValue + 1)
	}

	return nil
}

func GenerateEncodedLocbin(locFile types.LocFile, folder string) error {
	// Create out directory if it does not exist
	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		return err
	}

	pb := progress.CreateProgressBar(float64(len(locFile)), progress.ProgressConfig{Message: "Генерация .locbin файлов"})
	pb.Start()

	for _, fd := range locFile {
		if err := generateLocbinFile(fd, folder); err != nil {
			return err
		}

		pb.Update(pb.CurrentValue + 1)
	}

	return nil
}

func generateTxtFile(fileData models.FileData, outFolder string) error {
	filePath := filepath.Join(outFolder, fileData.Location)
	filePath = strings.Replace(filePath, ".locbin", ".txt", 1)

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

func generateLocbinFile(fileData models.FileData, outFolder string) error {
	filePath := filepath.Join(outFolder, fileData.Location)

	err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		return fmt.Errorf("Ошибка при создании папки %s: %v", filepath.Dir(filePath), err)
	}

	outFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("Ошибка при создании файла %s: %v", filePath, err)
	}
	defer outFile.Close()

	var memoryStream2 bytes.Buffer

	for _, kv := range fileData.Dictionary {
		key := kv.Key
		text := kv.Loc.Ru

		if text == "<##########>" {
			bytesData := []byte(key)
			var memoryStream bytes.Buffer
			memoryStream.WriteByte(10)
			if err := writeLength(&memoryStream, len(bytesData)); err != nil {
				return fmt.Errorf("Ошибка записи длины: %v", err)
			}
			memoryStream.Write(bytesData)

			array2 := memoryStream.Bytes()
			memoryStream2.WriteByte(10)
			if err := writeLength(&memoryStream2, len(array2)); err != nil {
				return fmt.Errorf("Ошибка записи длины: %v", err)
			}
			memoryStream2.Write(array2)
			continue
		}

		text = strings.ReplaceAll(text, "<lw>", "\n")
		text = strings.ReplaceAll(text, "<lwr>", "\r\n")

		bytes2 := []byte(key)
		bytes3 := []byte(text)

		var memoryStream3 bytes.Buffer
		memoryStream3.WriteByte(10)
		if err := writeLength(&memoryStream3, len(bytes2)); err != nil {
			return fmt.Errorf("Ошибка записи длины: %v", err)
		}
		memoryStream3.Write(bytes2)
		memoryStream3.WriteByte(18)
		if err := writeLength(&memoryStream3, len(bytes3)); err != nil {
			return fmt.Errorf("Ошибка записи длины: %v", err)
		}
		memoryStream3.Write(bytes3)

		array3 := memoryStream3.Bytes()
		memoryStream2.WriteByte(10)
		if err := writeLength(&memoryStream2, len(array3)); err != nil {
			return fmt.Errorf("Ошибка записи длины: %v", err)
		}
		memoryStream2.Write(array3)
	}

	if err := os.WriteFile(filePath, memoryStream2.Bytes(), os.ModePerm); err != nil {
		return fmt.Errorf("Ошибка записи файла locbin %s: %v", filePath, err)
	}

	return nil
}

func writeLength(stream *bytes.Buffer, length int) error {
	for length >= 128 {
		if err := stream.WriteByte(byte((length & 0x7F) | 0x80)); err != nil {
			return err
		}
		length >>= 7
	}
	return stream.WriteByte(byte(length))
}
