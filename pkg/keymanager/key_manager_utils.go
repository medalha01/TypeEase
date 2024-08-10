    package keymanager

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"github.com/go-vgo/robotgo"

	"fyne.io/fyne/v2/dialog"
)

// saveKeyEntry saves a key entry, either by creating a new one or updating an existing one.
func (km *KeyManager) saveKeyEntry() {
	keyName := km.nameEntry.Text
	keyDescription := km.descEntry.Text

	if keyName == "" {
		dialog.ShowInformation("Error", "Activation Key cannot be empty", km.window)
		return
	}

	if strings.Contains(keyName, " ") {
		dialog.ShowInformation("Error", "Activation Key cannot contain spaces", km.window)
		return
	}

	if strings.ContainsAny(keyName, "*(:") || strings.ContainsAny(keyDescription, "*):") {
		dialog.ShowInformation("Error", "Key Name and Description cannot contain '*(:' or '*):'", km.window)
		return
	}

	if existingEntry, exists := km.entries[keyName]; exists {
		dialog.ShowConfirm("Update Key", fmt.Sprintf("The key '%s' already exists. Do you want to update it?", keyName), func(confirm bool) {
			if confirm {
				existingEntry.Description = keyDescription
				existingEntry.Active = km.activeToggle.Checked
				km.refreshKeyList()
				km.saveEntriesToFile()
				km.clearInputs()
			}
		}, km.window)
	} else {
		km.createNewEntry(keyName, keyDescription)
		km.saveEntriesToFile()
		km.clearInputs()
	}
}

// loadEntriesFromFile loads key entries from the storage file.
func (km *KeyManager) LoadEntriesFromFile() error {
	file, err := os.Open(km.filePath)
	if err != nil {
		fmt.Printf("Failed to open key file: %v\n", err)
		if os.IsNotExist(err) {
			file, err = os.Create(km.filePath)
			if err != nil {
				return err
			}
			fmt.Printf("Created new key file: %s\n", km.filePath)
		} else {
			return err
		}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "*(:")
		if len(parts) == 2 {
			keyName := parts[0]
			descAndStatus := strings.Split(parts[1], "*):")
			if len(descAndStatus) == 2 {
				keyDescription := descAndStatus[0]
				activeStatus := strings.TrimSpace(descAndStatus[1])
				active := activeStatus == "1"
				km.entries[keyName] = &KeyEntry{Name: keyName, Description: keyDescription, Active: active}
				km.entryNames = append(km.entryNames, keyName)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading key file: %v", err)
	}

	km.refreshKeyList()
	return nil
}

// HandleInput processes the input string and simulates typing using robotgo.
func HandleInput(input string, km *KeyManager) {
	if entry, exists := km.entries[input]; exists && entry.Active {
		robotgo.TypeStr(entry.Description)
	}
}

