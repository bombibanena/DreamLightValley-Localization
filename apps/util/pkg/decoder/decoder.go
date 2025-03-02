package decoder

import (
	"fmt"

	"ddv_loc/pkg/generator"
	"ddv_loc/pkg/reader"
)

func Decode(format string, inFolder string, outPath string) error {
	// Read encoded files
	locFile, err := reader.ReadEncodedFromLocbin(inFolder)
	if err != nil {
		return err
	}

	// Generate output file
	switch format {
	case "json":
		err = generator.GenerateDecodedJSON(locFile, outPath)
	case "csv":
		err = generator.GenerateDecodedCSV(locFile, outPath)
	case "raw":
		err = generator.GenerateDecodedTxt(locFile, outPath)
	default:
		return fmt.Errorf("Неизвестный формат: %v\n", format)
	}

	if err != nil {
		return err
	}

	return nil
}
