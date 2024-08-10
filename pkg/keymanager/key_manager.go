package keymanager

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// KeyEntry represents a key entry saved in the KeyManager.
type KeyEntry struct {
	Name        string // The name of the key (activation key)
	Description string // The text associated with the key
	Active      bool   // Whether the key is active
}

// KeyManager handles the storage and management of key entries.
type KeyManager struct {
	entries       map[string]*KeyEntry // Maps key names to their entries
	entryNames    []string             // List of key names, used for display
	keyList       *widget.List         // Fyne widget displaying the list of keys
	nameEntry     *widget.Entry        // Input field for key name
	descEntry     *widget.Entry        // Input field for key description
	activeToggle  *widget.Check        // Toggle for key activation status
	saveButton    *widget.Button       // Button to save a key entry
	deleteButton  *widget.Button       // Button to delete a key entry
	window        fyne.Window          // Reference to the Fyne window
	filePath      string               // Path to the file where keys are stored
}

// NewKeyManager initializes a new KeyManager instance.
func NewKeyManager(w fyne.Window, filePath string) *KeyManager {
	km := &KeyManager{
		entries:      make(map[string]*KeyEntry),
		entryNames:   []string{},
		window:       w,
		nameEntry:    widget.NewEntry(),
		descEntry:    widget.NewMultiLineEntry(),
		activeToggle: widget.NewCheck("Active", nil),
		filePath:     filePath,
	}

	km.keyList = km.createKeyList()
	km.saveButton = km.createSaveButton()
	km.deleteButton = km.createDeleteButton()

	return km
}

func (km *KeyManager) Entries() map[string]*KeyEntry {
    return km.entries
}

