package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

// Converts the given event into an indented and human-readable JSON string.
func convertEventToJsonString(event NostrEvent) string {
	eventJson, err := json.Marshal(event)
	if err != nil {
		fmt.Println("Error converting event to JSON:", err)
		return ""
	}
	var eventJsonMap map[string]interface{}
	if err := json.Unmarshal(eventJson, &eventJsonMap); err != nil {
		fmt.Println("Error decoding event JSON:", err)
		return ""
	}
	eventJson, err = json.MarshalIndent(eventJsonMap, "", "\t")
	if err != nil {
		fmt.Println("Error converting event to JSON:", err)
		return ""
	}
	return string(eventJson)
}

// Saves the given event to a JSON file inside an "archive" folder in the same directory
// as the executable. The file is saved in a subfolder named after the event kind. And the
// file name is the event id with a ".json" extension.
func SaveEventToArchive(event NostrEvent) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("Error getting current file path.")
		return
	}
	thisFileDir := filepath.Dir(filename)

	jsonFilePath := filepath.Join(thisFileDir, ArchivesFolder, event.Pubkey, strconv.Itoa(event.Kind), event.Id+".json")
	if _, err := os.Stat(jsonFilePath); err == nil {
		// This means the file already exists.
		fmt.Println("File already exists:", jsonFilePath)
		return
	}
	// Create all the directories if they don't exist.
	err := os.MkdirAll(filepath.Dir(jsonFilePath), os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directories:", err)
		return
	}

	file, err := os.OpenFile(jsonFilePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	eventJson := convertEventToJsonString(event)
	_, err = fmt.Fprintln(file, eventJson)
	if err != nil {
		fmt.Println("Error writing event to file:", err)
		return
	}
}
