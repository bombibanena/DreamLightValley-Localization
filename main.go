package main

import (
	"flag"
	"fmt"
	"os"

	"ddv_loc/decoder"
	"ddv_loc/encoder"
)

const (
	version = "1.0.0" // Version
)

func main() {
	// Set flags
	mode := flag.String("mode", "decode", "Режим работы (decode или encode)")
	format := flag.String("format", "", "Формат вывода (json, csv или raw)")
	inPath := flag.String("in", "", "Входная папка")
	outPath := flag.String("out", "", "Выходная папка")
	showHelp := flag.Bool("help", false, "Показать справку")
	showVersion := flag.Bool("version", false, "Показать версию программы")

	// Parse arguments
	flag.Parse()

	// Check --help
	if *showHelp {
		printHelp()
		os.Exit(0)
	}

	// Check --version
	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	// Check required arguments
	if *format == "" || *inPath == "" || *outPath == "" {
		fmt.Println("Необходимо указать все аргументы: --mode {decode|encode} --format {json|csv} --in in_folder --out out_folder")
		fmt.Println("Используйте --help для получения справки.")
		os.Exit(1)
	}

	// Check values for --mode
	if *mode != "decode" && *mode != "encode" {
		fmt.Println("Недопустимое значение для --mode. Допустимые значения: decode, encode")
		os.Exit(1)
	}

	// Check values for --format
	if *format != "json" && *format != "csv" && *format != "raw" {
		fmt.Println("Недопустимое значение для --format. Допустимые значения: json, csv")
		os.Exit(1)
	}

	switch *mode {
	case "decode":
		if err := decoder.Decode(*format, *inPath, *outPath); err != nil {
			fmt.Printf("Ошибка при декодировании: %v\n", err)
			os.Exit(1)
		}
	case "encode":
		if err := encoder.Encode(*format, *inPath, *outPath); err != nil {
			fmt.Printf("Ошибка при кодировании: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println("Обработка завершена успешно.")
}

func printHelp() {
	fmt.Println("Использование: .\\ddv_loc.exe --mode {decode|encode} --format {json|csv|raw} --in in_folder --out out_folder")
	fmt.Println()
	fmt.Println("Опции:")
	fmt.Println("  --mode\tРежим (decode или encode)")
	fmt.Println("  --format\tФормат (json, csv или raw)")
	fmt.Println("  --in\t\tВходная папка или файл")
	fmt.Println("  --out\t\tВыходная папка")
	fmt.Println("  --help\tПоказать справку")
	fmt.Println("  --version\tПоказать версию")
	fmt.Println()
	fmt.Println("Режим decode:")
	fmt.Println("  --in\t\tПуть до папки с .locbin файлами")
	fmt.Println("  --out\t\tРезультат работы в виде loc.json, loc.csv или папка с .txt файлами ")
	fmt.Println()
	fmt.Println("Режим encode:")
	fmt.Println("  --in\t\tПуть до json, csv или .txt файлов")
	fmt.Println("  --out\t\tРезультат работы в виде .locbin файлов")
	fmt.Println()
	fmt.Println("Форматы:")
	fmt.Println("  json\t\tJSON - [{\"location\":\"/folder/file.locbin\",\"dictionary\":[{\"key\":\"key_1\",\"loc\":{\"en\":\"Value\",\"ru\":\"Значение\"}}]}]")
	fmt.Println("  csv\t\tCSV - location,key,en,ru")
	fmt.Println("  raw\t\tДекодированные .locbin файлы в .txt формате, с сохранением исходной структуры папок")
	fmt.Println()
	fmt.Println("Примеры")
	fmt.Println(" - json:")
	fmt.Println("decode: .\\ddv_loc.exe --mode decode --format json --in in_folder --out out_folder")
	fmt.Println("encode: .\\ddv_loc.exe --mode encode --format json --in in_folder/loc.json --out out_folder")
	fmt.Println()
	fmt.Println(" - csv:")
	fmt.Println("decode: .\\ddv_loc.exe --mode decode --format csv --in in_folder --out out_folder")
	fmt.Println("encode: .\\ddv_loc.exe --mode encode --format csv --in in_folder/loc.csv --out out_folder")
	fmt.Println()
	fmt.Println(" - raw:")
	fmt.Println("decode: .\\ddv_loc.exe --mode decode --format raw --in in_folder --out out_folder")
	fmt.Println("encode: .\\ddv_loc.exe --mode encode --format raw --in in_folder --out out_folder")
}

func printVersion() {
	fmt.Printf("Версия: %s\n", version)
}
