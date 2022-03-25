package tui

import "github.com/rivo/tview"

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

	t.tabsLayout = flex
}
