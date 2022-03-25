package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (t *Tui) setTabsLayout() {
	flex := tview.NewFlex().SetDirection(tview.FlexColumn)

	flex.AddItem(createButton("Info", func() {
		t.showInfo()
	}), 0, 1, false)
	flex.AddItem(createButton("Stdout Logs", func() {
		t.tailStdoutLogs()
	}), 0, 1, false)
	flex.AddItem(createButton("Stderr Logs", func() {
		t.tailStderrLogs()
	}), 0, 1, false)

	flex.SetBorder(true)
	flex.SetTitle("Tabs")

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case '1':
				t.app.SetFocus(flex.GetItem(0))
				t.showInfo()
			case '2':
				t.app.SetFocus(flex.GetItem(1))
				t.tailStdoutLogs()
			case '3':
				t.app.SetFocus(flex.GetItem(2))
				t.tailStderrLogs()
			}
		}

		return event
	})

	t.tabsLayout = flex
}
