package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (t *Tui) setFooter() {
	footer := tview.NewTextView()

	footer.SetBorder(true)
	footer.SetTitle("Keyboard Shortcuts")
	footer.SetTitleColor(tcell.ColorOrange)

	footer.SetText(
		`App - q: quit; Tab: navigate through panels;
Tabs panel - 1: Info; 2: Stdout Logs; 3: Stderr Logs`,
	)
	footer.SetTextColor(tcell.ColorOrange)
	footer.SetTextAlign(tview.AlignCenter)

	t.footerLayout = footer
}
