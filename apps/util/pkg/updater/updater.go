package updater

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"ddv_loc/pkg/app"
	"ddv_loc/pkg/generator"
	"ddv_loc/pkg/models"
	"ddv_loc/pkg/reader"
	"ddv_loc/pkg/translator"
	"ddv_loc/pkg/types"
)

func CheckUpdates(format string, inPathOld string, inPathNew string, outPath string) (bool, error) {
	var locFileOld types.LocFile
	var locFileNew types.LocFile
	var err error

	switch format {
	case "json":
		locFileOld, err = reader.ReadDecodedFromJSON(inPathOld)
		locFileNew, err = reader.ReadDecodedFromJSON(inPathNew)
	case "csv":
		locFileOld, err = reader.ReadDecodedFromCSV(inPathOld)
		locFileNew, err = reader.ReadDecodedFromCSV(inPathNew)
	case "raw":
		locFileOld, err = reader.ReadDecodedFromTxt(inPathOld)
		locFileNew, err = reader.ReadDecodedFromTxt(inPathNew)
	default:
		return false, fmt.Errorf("Неизвестный формат: %v\n", format)
	}

	if err != nil {
		return false, err
	}

	updates := compareFileData(locFileOld, locFileNew)
	if !updates.Any() {
		return false, nil
	}

	jsonData, err := json.MarshalIndent(updates, "", "  ")
	if err != nil {
		return false, fmt.Errorf("Error marshaling JSON: %v", err)
	}

	file, err := os.Create(outPath)
	if err != nil {
		return false, fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		return false, fmt.Errorf("error writing to file: %v", err)
	}

	return true, nil
}

func Patch(format string, inPath string, updatesFilePath string, outFolder string, translate bool) error {
	var locFile types.LocFile
	var updates *types.LocFileUpdates
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

	updates, err = reader.ReadUpdatesFromJSON(updatesFilePath)
	if err != nil {
		return err
	}

	if translate {
		translatedUpdates, err := translateUpdates(*updates)
		if err != nil {
			return err
		}

		updates = translatedUpdates
	}

	applyUpdates(&locFile, updates)

	switch format {
	case "json":
		err = generator.GenerateDecodedJSON(locFile, outFolder)
	case "csv":
		err = generator.GenerateDecodedCSV(locFile, outFolder)
	case "raw":
		err = generator.GenerateDecodedTxt(locFile, outFolder)
	default:
		return fmt.Errorf("Неизвестный формат: %v\n", format)
	}

	if err != nil {
		return err
	}

	return nil
}

func compareFileData(old types.LocFile, new types.LocFile) types.LocFileUpdates {
	var updates types.LocFileUpdates

	oldMap := make(map[string]models.FileData)
	for _, fd := range old {
		oldMap[fd.Location] = fd
	}

	newMap := make(map[string]models.FileData)
	for _, fd := range new {
		newMap[fd.Location] = fd
	}

	for loc, newFD := range newMap {
		if _, exists := oldMap[loc]; !exists {
			updates.New = append(updates.New, newFD)
		}
	}

	for loc, oldFD := range oldMap {
		if newFD, exists := newMap[loc]; exists {
			compareDictionaries(&updates, loc, oldFD.Dictionary, newFD.Dictionary)
		} else {
			updates.Removed = append(updates.Removed, oldFD)
		}
	}

	sortLocFile(&updates.New)
	sortLocFile(&updates.Changes.New)
	sortLocFile(&updates.Changes.Changes)
	sortLocFile(&updates.Changes.Removed)
	sortLocFile(&updates.Removed)

	return updates
}

func compareDictionaries(updates *types.LocFileUpdates, loc string, oldDict, newDict []models.KeyValue) {
	var newKeys, changedKeys, removedKeys []models.KeyValue

	oldMap := make(map[string]models.KeyValue)
	for _, kv := range oldDict {
		oldMap[kv.Key] = kv
	}

	newMap := make(map[string]models.KeyValue)
	for _, kv := range newDict {
		newMap[kv.Key] = kv
	}

	for key, newKV := range newMap {
		if oldKV, exists := oldMap[key]; exists {
			if oldKV.Loc.En != newKV.Loc.En {
				changedKeys = append(changedKeys, newKV)
			}
		} else {
			newKeys = append(newKeys, newKV)
		}
	}

	for key, oldKV := range oldMap {
		if _, exists := newMap[key]; !exists {
			removedKeys = append(removedKeys, oldKV)
		}
	}

	if len(newKeys) > 0 || len(changedKeys) > 0 || len(removedKeys) > 0 {
		fileData := models.FileData{
			Location: loc,
		}

		if len(newKeys) > 0 {
			fileData.Dictionary = newKeys
			updates.Changes.New = append(updates.Changes.New, fileData)
		}

		if len(changedKeys) > 0 {
			fileData.Dictionary = changedKeys
			updates.Changes.Changes = append(updates.Changes.Changes, fileData)
		}

		if len(removedKeys) > 0 {
			fileData.Dictionary = removedKeys
			updates.Changes.Removed = append(updates.Changes.Removed, fileData)
		}
	}
}

func sortLocFile(locFile *types.LocFile) {
	sort.Slice(*locFile, func(i, j int) bool {
		return (*locFile)[i].Location < (*locFile)[j].Location
	})

	for i := range *locFile {
		sort.Slice((*locFile)[i].Dictionary, func(j, k int) bool {
			return (*locFile)[i].Dictionary[j].Key < (*locFile)[i].Dictionary[k].Key
		})
	}
}

func applyUpdates(locFile *types.LocFile, updates *types.LocFileUpdates) {
	for _, newFD := range updates.New {
		*locFile = append(*locFile, newFD)
	}

	for _, changedFD := range updates.Changes.Changes {
		for i, fd := range *locFile {
			if fd.Location == changedFD.Location {
				(*locFile)[i].Dictionary = mergeDictionaries(fd.Dictionary, changedFD.Dictionary)
				break
			}
		}
	}

	for _, newKVs := range updates.Changes.New {
		for i, fd := range *locFile {
			if fd.Location == newKVs.Location {
				(*locFile)[i].Dictionary = mergeDictionaries(fd.Dictionary, newKVs.Dictionary)
				break
			}
		}
	}

	for _, removedKVs := range updates.Changes.Removed {
		for i, fd := range *locFile {
			if fd.Location == removedKVs.Location {
				(*locFile)[i].Dictionary = removeKeyValues(fd.Dictionary, removedKVs.Dictionary)
				break
			}
		}
	}

	for _, removedFD := range updates.Removed {
		for i := len(*locFile) - 1; i >= 0; i-- {
			if (*locFile)[i].Location == removedFD.Location {
				*locFile = append((*locFile)[:i], (*locFile)[i+1:]...)
				break
			}
		}
	}
}

func mergeDictionaries(oldDict, newDict []models.KeyValue) []models.KeyValue {
	result := make([]models.KeyValue, len(oldDict))
	copy(result, oldDict)

	newMap := make(map[string]models.KeyValue)
	for _, kv := range newDict {
		newMap[kv.Key] = kv
	}

	for i, kv := range result {
		if newKV, exists := newMap[kv.Key]; exists {
			result[i] = newKV
			delete(newMap, kv.Key)
		}
	}

	for _, kv := range newDict {
		if _, exists := newMap[kv.Key]; exists {
			result = append(result, kv)
		}
	}

	return result
}

func removeKeyValues(oldDict, toRemove []models.KeyValue) []models.KeyValue {
	removeMap := make(map[string]bool)
	for _, kv := range toRemove {
		removeMap[kv.Key] = true
	}

	result := []models.KeyValue{}
	for _, kv := range oldDict {
		if !removeMap[kv.Key] {
			result = append(result, kv)
		}
	}

	return result
}

func translateUpdates(updates types.LocFileUpdates) (*types.LocFileUpdates, error) {
	translatedUpdates := updates

	tr := translator.NewTranslator(app.Config.Translator, "RU")

	translatedUpdatesNew, err := tr.TranslateLocFile(translatedUpdates.New)
	if err != nil {
		return nil, err
	}

	translatedUpdatesChangesNew, err := tr.TranslateLocFile(translatedUpdates.Changes.New)
	if err != nil {
		return nil, err
	}

	translatedUpdatesChangesChanges, err := tr.TranslateLocFile(translatedUpdates.Changes.Changes)
	if err != nil {
		return nil, err
	}

	translatedUpdates.New = translatedUpdatesNew
	translatedUpdates.Changes.New = translatedUpdatesChangesNew
	translatedUpdates.Changes.Changes = translatedUpdatesChangesChanges

	return &translatedUpdates, nil
}
