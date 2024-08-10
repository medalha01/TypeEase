package keymanager

import (
	"bufio"
	"fmt"
	"os"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)


// createKeyList initializes the key list widget.
func (km *KeyManager) createKeyList() *widget.List {
	list := widget.NewList(
		func() int { return len(km.entryNames) },
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.ContentPasteIcon()), widget.NewLabel("Activation Key"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			km.updateKeyListItem(id, item)
		},
	)

	list.OnSelected = km.handleKeySelection

	return list
}

func removeKeyFromSlice(slice []string, key string) []string {
	for i, k := range slice {
		if k == key {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

// updateKeyListItem updates the display of a list item in the key list.
func (km *KeyManager) updateKeyListItem(id widget.ListItemID, item fyne.CanvasObject) {
	entry := km.entries[km.entryNames[id]]
	label := item.(*fyne.Container).Objects[1].(*widget.Label)
	label.SetText(entry.Name)

	if entry.Active {
		label.TextStyle = fyne.TextStyle{}
	} else {
		label.TextStyle = fyne.TextStyle{Italic: true}
	}
	label.Refresh()
}

// handleKeySelection loads the selected key entry's details into the input fields.
func (km *KeyManager) handleKeySelection(id widget.ListItemID) {
	selectedEntry := km.entries[km.entryNames[id]]
	if selectedEntry != nil {
		km.nameEntry.SetText(selectedEntry.Name)
		km.descEntry.SetText(selectedEntry.Description)
		km.activeToggle.SetChecked(selectedEntry.Active)
	}
}

// createSaveButton creates the Save button and defines its behavior.
func (km *KeyManager) createSaveButton() *widget.Button {
	return widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), km.saveKeyEntry)
}

// createDeleteButton creates the Delete button and defines its behavior.
func (km *KeyManager) createDeleteButton() *widget.Button {
	return widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), km.deleteKeyEntry)
}

// deleteKeyEntry deletes the selected key entry from the KeyManager.
func (km *KeyManager) deleteKeyEntry() {
	keyName := km.nameEntry.Text

	if keyName == "" {
		dialog.ShowInformation("Error", "No Activation Key selected", km.window)
		return
	}

	if _, exists := km.entries[keyName]; !exists {
		dialog.ShowInformation("Error", "The selected key does not exist", km.window)
		return
	}

	dialog.ShowConfirm("Delete Key", "Are you sure you want to delete this key?", func(confirm bool) {
		if confirm {
			delete(km.entries, keyName)
			km.entryNames = removeKeyFromSlice(km.entryNames, keyName)
			km.keyList.Refresh()
			km.clearInputs()
			km.saveEntriesToFile()
		}
	}, km.window)
}

// saveEntriesToFile writes all key entries to the storage file.
func (km *KeyManager) saveEntriesToFile() {
	file, err := os.Create(km.filePath)
	if err != nil {
		dialog.ShowError(err, km.window)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, keyName := range km.entryNames {
		entry := km.entries[keyName]
		activeStatus := 0
		if entry.Active {
			activeStatus = 1
		}
		_, err := fmt.Fprintf(writer, "%s*(:%s*):%d\n", entry.Name, entry.Description, activeStatus)
		if err != nil {
			dialog.ShowError(err, km.window)
			return
		}
	}
	writer.Flush()
}

// clearInputs clears the input fields in the GUI.
func (km *KeyManager) clearInputs() {
	km.nameEntry.SetText("")
	km.descEntry.SetText("")
	km.activeToggle.SetChecked(true)
}

// createNewEntry creates a new key entry and adds it to the KeyManager.
func (km *KeyManager) createNewEntry(keyName, keyDescription string) {
	km.entries[keyName] = &KeyEntry{Name: keyName, Description: keyDescription, Active: km.activeToggle.Checked}
	km.entryNames = append(km.entryNames, keyName)
	km.refreshKeyList()
}

// refreshKeyList refreshes the list of key entries displayed in the GUI.
func (km *KeyManager) refreshKeyList() {
	sort.Strings(km.entryNames)
	km.keyList.Refresh()
}


// CreateMainLayout constructs the main layout for the KeyManager application.
func (km *KeyManager) CreateMainLayout() fyne.CanvasObject {
	leftPane := container.NewBorder(
		widget.NewLabelWithStyle("Key List", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		nil, nil, nil, container.NewMax(container.NewVScroll(km.keyList)),
	)

	buttons := container.NewHBox(
		km.saveButton,
		km.deleteButton,
	)

	form := container.NewVBox(
		widget.NewLabelWithStyle("Key Details", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		widget.NewLabel("Activation Key:"),
		km.nameEntry,
		widget.NewLabel("Substitute text:"),
		km.descEntry,
		km.activeToggle,
		container.NewCenter(buttons),
	)

	formWithPadding := container.NewVBox(container.NewVBox(form))

	mainLayout := container.NewHSplit(leftPane, formWithPadding)
	mainLayout.SetOffset(0.4)

	return container.NewPadded(mainLayout)
}
