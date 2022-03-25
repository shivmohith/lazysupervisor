package tui

import "github.com/gdamore/tcell/v2"

func (t *Tui) captureKeyboardEvents() {
	t.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q':
				t.app.Stop()
			}
		case tcell.KeyTab:
			t.app.SetFocus(t.panels[t.currentPanelInFocusPosition])
			t.currentPanelInFocusPosition++

			if int(t.currentPanelInFocusPosition) > len(t.panels) {
				t.currentPanelInFocusPosition = 1
			}
		}

		return event
	})
}
