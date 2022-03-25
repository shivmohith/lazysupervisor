package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	infoTabPosition   = 0
	stdoutTabPosition = 1
	stderrTabPosition = 2
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
		if event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case '1':
				t.app.SetFocus(flex.GetItem(infoTabPosition))
				t.showInfo()
			case '2':
				t.app.SetFocus(flex.GetItem(stdoutTabPosition))
				t.tailStdoutLogs()
			case '3':
				t.app.SetFocus(flex.GetItem(stderrTabPosition))
				t.tailStderrLogs()
			}
		}

		return event
	})

	t.tabsLayout = flex
}
