package encoder

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ddv_loc/models"
)

func generateEncodedFiles(fileData []models.FileData, outPath string) error {
	for _, fd := range fileData {
		if err := generateEncodedFile(fd, outPath); err != nil {
			fmt.Printf("    Encode - Error\n")
			return fmt.Errorf("Ошибка при генерации файла %s: %v", fd, err)
		}

		fmt.Printf("    Encode - OK\n")
	}

	return nil
}

func generateEncodedFile(fileData models.FileData, outFolder string) error {
	filePath := filepath.Join(outFolder, fileData.Location)

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
