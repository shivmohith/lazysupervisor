package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (t *Tui) setHeader() {
	header := tview.NewTextView()
	header.SetText("Welcome to Supervisor TUI")
	header.SetTextAlign(tview.AlignCenter)
	header.SetBorderPadding(1, 1, 1, 1)
	header.SetBackgroundColor(tcell.ColorBlueViolet)

	t.headerLayout = header
}
