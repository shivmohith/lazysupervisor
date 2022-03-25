package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	log "github.com/sirupsen/logrus"
)

func (t *Tui) setProcessLayout() {
	list := tview.NewList()
	list.ShowSecondaryText(false)
	list.SetBorder(true)
	list.SetTitle("Process")

	if t.groupLayout.GetItemCount() > 0 {
		g, _ := t.groupLayout.GetItemText(t.groupLayout.GetCurrentItem())
		for _, p := range t.groupToProcessesInfo[g] {
			list.AddItem(p.Name, "", 0, nil)
		}
	}

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case 'm':
				t.processModal.SetText(fmt.Sprintf(`Group Selected: %s
Process Selected: %s`, t.selectedGroup, t.selectedProcess))

				pages := tview.NewPages().
					AddPage("background", t.layout, true, true).
					AddPage("modal", t.processModal, true, true)
				t.app.SetRoot(pages, true)
			}
		}

		return event
	})

	t.processLayout = list

	modal := tview.NewModal().
		AddButtons([]string{"Start Process", "Stop Process", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			processName := fmt.Sprintf("%s:%s", t.selectedGroup, t.selectedProcess)

			switch buttonLabel {
			case "Start Process":
				if err := t.client.StartProcess(processName, false); err != nil {
					log.Errorf("starting process %s: %v", processName, err)
				}
			case "Stop Process":
				if err := t.client.StopProcess(processName, false); err != nil {
					log.Errorf("stopping process %s: %v", processName, err)
				}
			}

			t.app.SetRoot(t.layout, true).SetFocus(t.processLayout)
		})

	t.processModal = modal

}

func (t *Tui) handleProcessSelect(process string) {
	t.selectedProcess = process
	t.showInfo()
}
