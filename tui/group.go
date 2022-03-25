package tui

import "github.com/rivo/tview"

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

	t.groupLayout = list
}

func (t *Tui) handleGroupSelect(group string) {
	t.selectedGroup = group
	t.setProcesses(group)
	t.selectedProcess, _ = t.processLayout.GetItemText(t.processLayout.GetCurrentItem())
	t.showInfo()
}
