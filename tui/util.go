package tui

import (
	"time"

	"github.com/rivo/tview"
)

func runEvery(d time.Duration, f func()) {
	go func() {
		for range time.Tick(d) {
			f()
		}
	}()
}

func createButton(name string, handler func()) *tview.Button {
	button := tview.NewButton(name).SetSelectedFunc(handler)
	button.SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)

	return button
}
