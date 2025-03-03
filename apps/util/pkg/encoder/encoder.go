package encoder

import (
	"fmt"

	"ddv_loc/pkg/generator"
	"ddv_loc/pkg/reader"
	"ddv_loc/pkg/types"
)

func Encode(format string, inPath string, outPath string) error {
	var locFile types.LocFile
	var err error

	switch format {
	case "json":
		locFile, err = reader.ReadDecodedFromJSON(inPath)
	case "csv":
		locFile, err = reader.ReadDecodedFromCSV(inPath)
	case "raw":
		locFile, err = reader.ReadDecodedFromTxt(inPath)
	default:
		return fmt.Errorf("Неизвестный формат: %v\n", format)
	}

	if err != nil {
		return err
	}

	if err := generator.GenerateEncodedLocbin(locFile, outPath); err != nil {
		return err
	}

	return nil
}
