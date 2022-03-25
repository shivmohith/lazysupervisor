package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	log "github.com/sirupsen/logrus"
)

func (t *Tui) setGroupLayout() {
	list := tview.NewList()
	list.ShowSecondaryText(false)
	list.SetBorder(true)
	list.SetTitle("Group")

	for g := range t.groupToProcessesInfo {
		gCopy := g
		list.AddItem(g, "", 0, func() {
			t.handleGroupSelect(gCopy)
		})
	}

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case 'm':
				t.groupModal.SetText(fmt.Sprintf("Group Selected: %s", t.selectedGroup))

				pages := tview.NewPages().
					AddPage("background", t.layout, true, true).
					AddPage("modal", t.groupModal, true, true)
				t.app.SetRoot(pages, true)
			}
		}

		return event
	})

	t.groupLayout = list

	modal := tview.NewModal().
		AddButtons([]string{"Start Group", "Stop Group", "Remove Group", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			switch buttonLabel {
			case "Start Group":
				if _, err := t.client.StartProcessGroup(t.selectedGroup, true); err != nil {
					log.Errorf("starting group %s: %v", t.selectedGroup, err)
				}
			case "Stop Group":
				if _, err := t.client.StopProcessGroup(t.selectedGroup, true); err != nil {
					log.Errorf("stopping group %s: %v", t.selectedGroup, err)
				}
			case "Remove Group":
				err := t.client.RemoveProcessGroup(t.selectedGroup)
				if err != nil {
					log.Errorf("removing group %s: %v", t.selectedGroup, err)
				} else {
					delete(t.groupToProcessesInfo, t.selectedGroup)

					indices := t.groupLayout.FindItems(t.selectedGroup, "", false, false)
					if len(indices) > 0 {
						t.groupLayout.RemoveItem(indices[0])
					}

					t.processLayout.Clear()
				}
			}

			t.app.SetRoot(t.layout, true).SetFocus(t.groupLayout)
		})

	t.groupModal = modal
}

func (t *Tui) handleGroupSelect(group string) {
	t.selectedGroup = group
	t.setProcesses(group)
	if t.processLayout.GetItemCount() > 0 {
		t.selectedProcess, _ = t.processLayout.GetItemText(t.processLayout.GetCurrentItem())
		t.showInfo()
	}
}

func (t *Tui) setProcesses(group string) {
	t.processLayout.Clear()

	for p := range t.groupToProcessesInfo[group] {
		pCopy := p
		t.processLayout.AddItem(pCopy, "", 0, func() {
			t.handleProcessSelect(pCopy)
		})
	}
}
