package updater

import (
	"fmt"
	"os"
	"sort"

	"ddv_loc/pkg/models"
	"ddv_loc/pkg/reader"
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

	report := compareFileData(locFileOld, locFileNew)
	if report == "" {
		return false, nil
	}

	if err := writeReportToFile(outPath, report); err != nil {
		return false, err
	}

	return true, nil
}

func compareFileData(locFileOld types.LocFile, locFileNew types.LocFile) string {
	report := ""

	oldMap := make(map[string]models.FileData)
	for _, fd := range locFileOld {
		oldMap[fd.Location] = fd
	}

	newMap := make(map[string]models.FileData)
	for _, fd := range locFileNew {
		newMap[fd.Location] = fd
	}

	var newLocations, changedLocations, removedLocations []string

	for loc := range newMap {
		if _, exists := oldMap[loc]; !exists {
			newLocations = append(newLocations, loc)
		}
	}

	for loc := range oldMap {
		if _, exists := newMap[loc]; !exists {
			removedLocations = append(removedLocations, loc)
		} else {
			changedLocations = append(changedLocations, loc)
		}
	}

	sort.Strings(newLocations)
	sort.Strings(changedLocations)
	sort.Strings(removedLocations)

	for _, loc := range newLocations {
		report += fmt.Sprintf("[NEW]     %s\n", loc)
		newFD := newMap[loc]
		sort.Slice(newFD.Dictionary, func(i, j int) bool {
			return newFD.Dictionary[i].Key < newFD.Dictionary[j].Key
		})
		for _, kv := range newFD.Dictionary {
			report += fmt.Sprintf("      [NEW]     %s\n", kv.Key)
			report += fmt.Sprintf("              - %s\n", kv.Loc.En)
		}
		report += "\n"
	}

	for _, loc := range changedLocations {
		report += compareDictionaries(loc, oldMap[loc].Dictionary, newMap[loc].Dictionary)
	}

	for _, loc := range removedLocations {
		report += fmt.Sprintf("[REMOVED] %s\n", loc)
		oldFD := oldMap[loc]
		sort.Slice(oldFD.Dictionary, func(i, j int) bool {
			return oldFD.Dictionary[i].Key < oldFD.Dictionary[j].Key
		})
		for _, kv := range oldFD.Dictionary {
			report += fmt.Sprintf("      [REMOVED] %s\n", kv.Key)
			report += fmt.Sprintf("              - %s\n", kv.Loc.En)
		}
		report += "\n"
	}

	return report
}

func compareDictionaries(loc string, oldDict []models.KeyValue, newDict []models.KeyValue) string {
	report := ""
	oldMap := make(map[string]models.KeyValue)
	for _, kv := range oldDict {
		oldMap[kv.Key] = kv
	}

	newMap := make(map[string]models.KeyValue)
	for _, kv := range newDict {
		newMap[kv.Key] = kv
	}

	var newKeys, changedKeys, removedKeys []string

	for key, newKV := range newMap {
		if oldKV, exists := oldMap[key]; exists {
			if oldKV.Loc.En != newKV.Loc.En || oldKV.Loc.Ru != newKV.Loc.Ru {
				changedKeys = append(changedKeys, key)
			}
		} else {
			newKeys = append(newKeys, key)
		}
	}

	for key := range oldMap {
		if _, exists := newMap[key]; !exists {
			removedKeys = append(removedKeys, key)
		}
	}

	sort.Strings(newKeys)
	sort.Strings(changedKeys)
	sort.Strings(removedKeys)

	if len(newKeys) > 0 || len(changedKeys) > 0 || len(removedKeys) > 0 {
		report += fmt.Sprintf("[CHANGED] %s\n", loc)
	}

	for _, key := range newKeys {
		newKV := newMap[key]
		report += fmt.Sprintf("      [NEW]     %s\n", key)
		report += fmt.Sprintf("              - %s\n", newKV.Loc.En)
	}

	for _, key := range changedKeys {
		oldKV := oldMap[key]
		newKV := newMap[key]
		report += fmt.Sprintf("      [CHANGED] %s\n", key)
		report += fmt.Sprintf("              - %s\n", oldKV.Loc.En)
		report += fmt.Sprintf("              - %s\n", newKV.Loc.En)
	}

	for _, key := range removedKeys {
		oldKV := oldMap[key]
		report += fmt.Sprintf("      [REMOVED] %s\n", key)
		report += fmt.Sprintf("              - %s\n", oldKV.Loc.En)
	}

	if len(newKeys) > 0 || len(changedKeys) > 0 || len(removedKeys) > 0 {
		report += "\n"
	}

	return report
}

func writeReportToFile(filePath string, report string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("Error creating file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(report)
	if err != nil {
		return fmt.Errorf("Error writing to file: %v", err)
	}

	return nil
}
