package translator

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"ddv_loc/pkg/app"
	"ddv_loc/pkg/generator"
	"ddv_loc/pkg/models"
	"ddv_loc/pkg/reader"
	"ddv_loc/pkg/translator/client"
	"ddv_loc/pkg/translator/config"
	"ddv_loc/pkg/types"
	"ddv_loc/pkg/utils/progress"
)

type (
	translator struct {
		opts   translateOpts
		client client.IClient
		cfg    config.Config
	}

	translateOpts struct {
		sourceLang string
		targetLang string
	}
)

func newTranslator(cfg config.Config, targetLang string) *translator {
	return &translator{
		opts: translateOpts{
			sourceLang: "EN",
			targetLang: targetLang,
		},
		client: client.GetClient(cfg),
		cfg:    cfg,
	}
}

func Translate(format string, inPath string, outPath string, language string) error {
	var locFile types.LocFile
	var err error

	// Read input file(s)
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

	translator := newTranslator(app.Config.Translator, language)

	translated, err := translator.translateLocFile(locFile)
	if err != nil {
		return err
	}

	// Generate output file(s)
	switch format {
	case "json":
		err = generator.GenerateDecodedJSON(translated, outPath)
	case "csv":
		err = generator.GenerateDecodedCSV(translated, outPath)
	case "raw":
		err = generator.GenerateDecodedTxt(translated, outPath)
	default:
		return fmt.Errorf("Неизвестный формат: %v\n", format)
	}

	if err != nil {
		return err
	}

	return nil
}

func (t *translator) translateLocFile(locFile types.LocFile) (types.LocFile, error) {
	// Take texts from locFile
	var texts []string
	for _, fd := range locFile {
		for _, kv := range fd.Dictionary {
			texts = append(texts, kv.Loc.En)
		}
	}

	// Translate texts
	translatedTexts, err := t.translateTexts(texts)
	if err != nil {
		return nil, err
	}

	// Fill result
	translatedIndex := 0
	translatedLocFile := make(types.LocFile, len(locFile))
	for i, fd := range locFile {
		translatedDict := make([]models.KeyValue, len(fd.Dictionary))
		for j, kv := range fd.Dictionary {
			translatedDict[j] = models.KeyValue{
				Key: kv.Key,
				Loc: models.KeyValueLoc{
					En: kv.Loc.En,
					Ru: translatedTexts[translatedIndex],
				},
			}
			translatedIndex++
		}

		translatedLocFile[i] = models.FileData{
			Location:   fd.Location,
			Dictionary: translatedDict,
		}
	}

	return translatedLocFile, nil
}

func (t *translator) translateTexts(texts []string) ([]string, error) {
	pb := progress.CreateProgressBar(float64(len(texts)), progress.ProgressConfig{Message: "Перевод"})
	pb.Start()

	var batches [][]string
	var currentBatch []string
	currentBatchSize := 0

	textSize := t.cfg.Translation.TextSize
	batchSize := t.cfg.Translation.BatchSize
	separator := t.cfg.Translation.Separator

	// Split on batches
	for _, text := range texts {
		if utf8.RuneCountInString(text) > textSize {
			// Split long text
			parts := splitText(text, textSize)
			for _, part := range parts {
				if currentBatchSize+utf8.RuneCountInString(part)+len(separator) > batchSize {
					batches = append(batches, currentBatch)
					currentBatch = []string{}
					currentBatchSize = 0
				}

				currentBatch = append(currentBatch, part)
				currentBatchSize += utf8.RuneCountInString(part) + len(separator)
			}
		} else {
			if currentBatchSize+utf8.RuneCountInString(text)+len(separator) > batchSize {
				batches = append(batches, currentBatch)
				currentBatch = []string{}
				currentBatchSize = 0
			}

			currentBatch = append(currentBatch, text)
			currentBatchSize += utf8.RuneCountInString(text) + len(separator)
		}
	}

	if len(currentBatch) > 0 {
		batches = append(batches, currentBatch)
	}

	// Translate batch
	var translatedTexts []string
	for _, batch := range batches {
		translatedBatch, err := t.translateBatch(batch)
		if err != nil {
			pb.Stop()
			return nil, err
		}

		translatedTexts = append(translatedTexts, translatedBatch...)

		pb.Update(float64(len(translatedTexts)))
	}

	return translatedTexts, nil
}

func (t *translator) translateBatch(batch []string) ([]string, error) {
	separator := t.cfg.Translation.Separator

	// Combine batch into single text with separator
	text := strings.Join(batch, separator)

	res, err := t.translate(text)
	if err != nil {
		return nil, err
	}

	// Split translated text by separator
	translatedTexts := strings.Split(res, separator)

	if len(translatedTexts) != len(batch) {
		mid := len(batch) / 2

		firstHalf, err := t.translateBatch(batch[:mid])
		if err != nil {
			return nil, err
		}

		secondHalf, err := t.translateBatch(batch[mid:])
		if err != nil {
			return nil, err
		}

		return append(firstHalf, secondHalf...), nil
	}

	return translatedTexts, nil
}

func (t *translator) translate(text string) (string, error) {
	exceptions := t.cfg.Translation.GetExceptions()

	htmlTagRegex := regexp.MustCompile(`(<[^>]+>)|(</[^>]+>)`) // HTML-like tags
	placeholderRegex := regexp.MustCompile(`\{[^}]+\}`)        // Placeholders
	// mdRegex := regexp.MustCompile(`\*[^}]+\*`)                 // Markdown-like placeholders

	placeholders := make(placeholders)
	counter := 0

	// Replace regexp
	replaceMatches := func(word string, regex *regexp.Regexp) string {
		return regex.ReplaceAllStringFunc(word, func(match string) string {
			placeholder := placeholders.add(counter, match)

			counter++

			return placeholder
		})
	}

	outText := text
	outText = replaceMatches(outText, htmlTagRegex)
	outText = replaceMatches(outText, placeholderRegex)
	// outText = replaceMatches(outText, mdRegex)

	// Replace exceptions
	for word := range exceptions {
		placeholder := placeholders.add(counter, exceptions[word])

		outText = strings.ReplaceAll(outText, word, placeholder)

		counter++
	}

	// Escape special chars
	outText = escapeSpecialChars(outText)

	res, err := t.client.Translate(outText, t.opts.sourceLang, t.opts.targetLang)
	if err != nil {
		return "", err
	}

	// Unescape special chars
	resText := res
	resText = normalizePlaceholders(resText)
	resText = unescapeSpecialChars(resText)

	for placeholder, translation := range placeholders {
		resText = strings.ReplaceAll(resText, placeholder, translation)
	}

	return resText, nil
}
