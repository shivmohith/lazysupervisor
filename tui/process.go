package tui

import "github.com/rivo/tview"

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

	t.processLayout = list
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

func (t *Tui) handleProcessSelect(process string) {
	t.selectedProcess = process
	t.showInfo()
}
