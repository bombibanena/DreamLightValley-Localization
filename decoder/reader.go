package decoder

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/charmap"

	"ddv_loc/models"
)

func readEncodedFiles(folder string) ([]models.FileData, error) {
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

		// Skip not .locbin files
		if filepath.Ext(info.Name()) != ".locbin" {
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
		keyValues, err := parseEncodedFile(data)
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

func parseEncodedFile(data []byte) ([]models.KeyValue, error) {
	defer func() {
		os.RemoveAll("_tmp")
	}()

	// Temp block dir
	blocksDir := filepath.Join("_tmp", "blocks")
	if err := os.MkdirAll(blocksDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("Ошибка при создании папки %s: %v", blocksDir, err)
	}

	// Temp loc file
	locFileName := filepath.Join("_tmp", "loc")

	if err := os.MkdirAll(filepath.Dir(locFileName), os.ModePerm); err != nil {
		return nil, fmt.Errorf("Ошибка при создании папки %s: %v", filepath.Dir(locFileName), err)
	}

	locFile, err := os.Create(locFileName)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при создании файла %s: %v", locFileName, err)
	}
	defer locFile.Close()

	numBlocks := 0
	for i := 0; i < len(data); {
		if data[i] == 10 {
			i++
			if i >= len(data) {
				break
			}

			num3, num4 := 0, 1
			for ; i < len(data); i++ {
				b := data[i]
				num3 += int(b&0x7F) * num4
				if b&0x80 == 0 {
					i++
					break
				}
				num4 *= 128
			}

			if i+num3 > len(data) {
				break
			}

			blockData := make([]byte, num3)
			copy(blockData, data[i:i+num3])

			blockFileName := fmt.Sprintf("block_%d.bin", numBlocks+1)
			blockFilePath := filepath.Join(blocksDir, blockFileName)

			if err := os.WriteFile(blockFilePath, blockData, os.ModePerm); err != nil {
				return nil, fmt.Errorf("Ошибка при записи блока %s: %v", blockFilePath, err)
			}

			blockTextFileName := fmt.Sprintf("block_%d.txt", numBlocks+1)
			blockTextFilePath := filepath.Join(blocksDir, blockTextFileName)

			var text3 string
			if num3 > 128 {
				num5 := int(math.Ceil(float64(num3) / 128.0))
				text3 = fmt.Sprintf("%X", num3+num5*128)
			} else {
				text3 = fmt.Sprintf("%X", num3)
			}

			if len(text3)%2 != 0 {
				text3 = "0" + text3
			}

			num6 := len(text3) / 2
			if err := os.WriteFile(blockTextFilePath, []byte(fmt.Sprintf("%d", num6)), os.ModePerm); err != nil {
				return nil, fmt.Errorf("Ошибка при записи блока в файл %s: %v", blockTextFilePath, err)
			}

			if err := processBlock(blockData, num6, locFile); err != nil {
				return nil, fmt.Errorf("Ошибка обработки блока: %v", err)
			}

			i += num3
			numBlocks++
		} else {
			i++
		}
	}

	// Read file
	content, err := os.ReadFile(locFileName)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при чтении файла %s: %v", locFileName, err)
	}

	var keyValues []models.KeyValue

	// Parse file
	lines := strings.Split(string(content), "\n")

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

func processBlock(block []byte, byteCount int, writer *os.File) error {
	num := 0
	iso88591 := charmap.ISO8859_1.NewDecoder()

	for num < len(block) {
		if block[num] == 10 {
			num++
			if num >= len(block) {
				break
			}

			num2 := int(block[num])
			num++
			if num+num2 > len(block) {
				break
			}

			str, err := iso88591.Bytes(block[num : num+num2])
			if err != nil {
				return fmt.Errorf("Ошибка чтения строки: %v", err)
			}
			num += num2

			if num >= len(block) {
				writer.WriteString("TAG: " + string(str) + "\n")
				writer.WriteString("<##########>\n\n")
				break
			}

			_ = block[num]
			num++

			num3 := num
			num += byteCount

			if byteCount > 1 && (num >= len(block) || len(block)-num < 128) {
				num = num3 + (byteCount - 1)
				if num >= len(block) {
					break
				}
			}

			str2, err := iso88591.Bytes(block[num:])
			if err != nil {
				return fmt.Errorf("Ошибка чтения строки: %v", err)
			}
			num = len(block)

			str2Replaced := strings.NewReplacer(
				"\r\n", "<lwr>",
				"\n", "<lw>",
				"\u0001", "",
			).Replace(string(str2))

			writer.WriteString("TAG: " + string(str) + "\n")
			if str2Replaced != "" {
				writer.WriteString(str2Replaced + "\n")
			} else {
				writer.WriteString("<##########>\n")
			}
			writer.WriteString("\n")
		} else {
			num++
		}
	}

	return nil
}
