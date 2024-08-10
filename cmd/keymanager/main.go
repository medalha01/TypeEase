package main

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"github.com/medalha01/keymanager/pkg/keymanager"
	"github.com/medalha01/keymanager/pkg/listener"
)

func main() {
	devices := []string{
		"/dev/input/event0",
		"/dev/input/event1",
		"/dev/input/event2",
		"/dev/input/event3",
		"/dev/input/event4",
		"/dev/input/event5",
		"/dev/input/event6",
		"/dev/input/event7",
	}

	var wg sync.WaitGroup
	inputChannel := make(chan string)

	a := app.New()
	w := a.NewWindow("Key Manager")
	a.Settings().SetTheme(theme.DarkTheme())

	keyManager := keymanager.NewKeyManager(w, "data.enf")

	if err := keyManager.LoadEntriesFromFile(); err != nil {
		dialog.ShowError(err, w)
	}

	w.SetContent(keyManager.CreateMainLayout())
	w.Resize(fyne.NewSize(700, 450))
	w.CenterOnScreen()

	for i, device := range devices {
		wg.Add(1)
		go listener.ListenToDevice(device, i, &wg, inputChannel, keyManager)
	}

	go func() {
		for input := range inputChannel {
			keymanager.HandleInput(input, keyManager)
		}
	}()

	w.ShowAndRun()

	wg.Wait()
	close(inputChannel)
}

