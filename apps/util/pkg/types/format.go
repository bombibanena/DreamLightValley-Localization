package types

import (
	"errors"
	"fmt"
	"strings"
)

type FormatEnum string

const (
	json FormatEnum = "json"
	csv  FormatEnum = "csv"
	raw  FormatEnum = "raw"
)

func (e *FormatEnum) String() string {
	return string(*e)
}

func (e *FormatEnum) Set(v string) error {
	switch v {
	case string(json), string(csv), string(raw):
		*e = FormatEnum(v)
		return nil
	default:
		return errors.New(fmt.Sprintf(`must be one of: "%s", "%s" or "%s"`, json, csv, raw))
	}
}

func (e *FormatEnum) Type() string {
	return "{" + strings.Join([]string{string(json), string(csv), string(raw)}, "|") + "}"
}
