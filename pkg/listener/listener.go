package listener

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"sync"

	"github.com/go-vgo/robotgo"
	"github.com/medalha01/keymanager/pkg/keymanager"
)

// InputEvent represents an input event from /dev/input devices.
type InputEvent struct {
	Timestamp [16]byte // Placeholder for timestamp (64-bit systems)
	Type      uint16   // Event type
	Code      uint16   // Key code or event code
	Value     int32    // Value associated with the event (e.g., key press/release)
}

// KeyMap maps key codes to their respective key names.
var KeyMap = map[uint16]string{
	1: "ESC", 2: "1", 3: "2", 4: "3", 5: "4", 6: "5", 7: "6", 8: "7", 9: "8",
	10: "9", 11: "0", 12: "-", 13: "=", 14: "BACKSPACE", 15: "TAB",
	16: "q", 17: "w", 18: "e", 19: "r", 20: "t", 21: "y", 22: "u", 23: "i",
	24: "o", 25: "p", 26: "[", 27: "]", 28: "ENTER", 29: "CTRL",
	30: "a", 31: "s", 32: "d", 33: "f", 34: "g", 35: "h", 36: "j", 37: "k",
	38: "l", 39: ";", 40: "'", 41: "`", 42: "SHIFT", 43: "\\",
	44: "z", 45: "x", 46: "c", 47: "v", 48: "b", 49: "n", 50: "m",
	51: ",", 52: ".", 53: "/", 54: "SHIFT", 55: "*", 56: "ALT",
	57: "SPACE", 58: "CAPSLOCK",
}

// ListenToDevice listens for input events from a device and processes them.
func ListenToDevice(device string, deviceID int, wg *sync.WaitGroup, inputChannel chan<- string, km *keymanager.KeyManager) {
	defer wg.Done()

	file, err := os.Open(device)
	if err != nil {
		fmt.Printf("Error opening input device %s: %v\n", device, err)
		return
	}
	defer file.Close()

	event := InputEvent{}
	eventSize := binary.Size(event)
	buf := make([]byte, eventSize)

	fmt.Printf("Listening on device %d: %s\n", deviceID, device)

	var currentInput string

	for {
		_, err := file.Read(buf)
		if err != nil {
			fmt.Printf("Error reading input event from %s: %v\n", device, err)
			break
		}

		reader := bytes.NewReader(buf)
		err = binary.Read(reader, binary.LittleEndian, &event)
		if err != nil {
			fmt.Printf("Error decoding input event from %s: %v\n", device, err)
			break
		}

		if event.Type == 1 && event.Value == 1 {
			key := KeyMap[event.Code]
			if key == "SHIFT" || key == "CTRL" || key == "ALT" || key == "CAPSLOCK" || key == "TAB" {
				// Ignore these keys
				fmt.Printf("Device %d: Modifier key pressed: %s\n", deviceID, key)
			} else if key != "" {
				fmt.Printf("Device %d: Key pressed: %s\n", deviceID, key)

				if key == "SPACE" || key == "ENTER" {
					if entry, exists := km.Entries()[currentInput]; exists && entry.Active {
						inputLength := len(currentInput)

						for i := 0; i < inputLength+1; i++ {
							robotgo.KeyTap("backspace")
						}
						robotgo.TypeStr(entry.Description)
					}
					currentInput = ""
				} else if key == "BACKSPACE" {
					if len(currentInput) > 0 {
						currentInput = currentInput[:len(currentInput)-1]
					}
				} else {
					currentInput += key
				}
			} else {
				fmt.Printf("Device %d: Unknown key code: %d\n", deviceID, event.Code)
			}
		}
	}
}

