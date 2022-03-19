package tui

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	gosupervisord "github.com/shivmohith/go-supervisord"
	"github.com/shivmohith/tui-supervisor/supervisord"
	log "github.com/sirupsen/logrus"
)

type tui struct {
	client supervisord.Client

	groupToProcessesInfo map[string]map[string]gosupervisord.ProcessInfo

	app    *tview.Application
	layout *tview.Flex

	headerLayout *tview.TextView
	footerLayout *tview.TextView

	groupLayout   *tview.List
	processLayout *tview.List
	infoTextView  *tview.TextView

	selectedGroup   string
	selectedProcess string
}

func New(client supervisord.Client) *tui {
	t := new(tui)

	t.client = client

	t.app = tview.NewApplication()

	t.setHeader()

	t.groupToProcessesInfo = t.getGroupProcessesMap()
	t.setGroupLayout()
	t.setProcessLayout()
	t.setInfoTextView()

	t.setLayout()

	t.selectedGroup, _ = t.groupLayout.GetItemText(t.groupLayout.GetCurrentItem())
	t.selectedProcess, _ = t.processLayout.GetItemText(t.processLayout.GetCurrentItem())
	return t
}

func (t *tui) refreshGroupsAndProcesses() {
	gTop := t.getGroupProcessesMap()

	// Adds new groups created to the group list
	for g, pMap := range gTop {
		gCopy := g
		if _, ok := t.groupToProcessesInfo[gCopy]; !ok {
			t.groupToProcessesInfo[gCopy] = pMap
			t.groupLayout.AddItem(gCopy, "", 0, func() {
				t.handleGroupSelect(gCopy)
			})
		}
	}

	// Remove groups from the group list that have removed from supervisord
	for g, _ := range t.groupToProcessesInfo {
		gCopy := g
		if _, ok := gTop[gCopy]; !ok {
			indices := t.groupLayout.FindItems(gCopy, "", false, false)
			if len(indices) == 0 {
				continue
			}

			t.groupLayout.RemoveItem(indices[0])
		}
	}

}

func (t *tui) Start() error {
	log.Info("Starting the tui app")

	runEvery(1*time.Second, func() {
		t.refreshGroupsAndProcesses()
	})

	return t.app.SetRoot(t.layout, true).SetFocus(t.layout).EnableMouse(true).Run()
}

func (t *tui) setHeader() {
	header := tview.NewTextView()
	header.SetText("Welcome to Supervisor TUI")
	header.SetTextAlign(tview.AlignCenter)
	header.SetBorderPadding(1, 1, 1, 1)
	header.SetBackgroundColor(tcell.ColorBlueViolet)

	t.headerLayout = header
}

func (t *tui) setLayout() {
	t.layout = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(t.headerLayout, 0, 1, false).
		AddItem(t.getSupervisorLayout(), 0, 12, false)
}

func (t *tui) setGroupLayout() {
	list := tview.NewList()
	list.ShowSecondaryText(false)
	list.SetBorder(true)
	list.SetTitle("Group")

	for g, _ := range t.groupToProcessesInfo {
		gCopy := g
		list.AddItem(g, "", 0, func() {
			t.handleGroupSelect(gCopy)
		})
	}

	t.groupLayout = list
}

func (t *tui) handleGroupSelect(group string) {
	t.selectedGroup = group
	t.setProcesses(group)
	t.selectedProcess, _ = t.processLayout.GetItemText(t.processLayout.GetCurrentItem())
	t.showInfo()
}

func (t *tui) setProcessLayout() {
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

func (t *tui) setProcesses(group string) {
	t.processLayout.Clear()

	for p, _ := range t.groupToProcessesInfo[group] {
		pCopy := p
		t.processLayout.AddItem(pCopy, "", 0, func() {
			t.handleProcessSelect(pCopy)
		})
	}
}

func (t *tui) handleProcessSelect(process string) {
	t.selectedProcess = process
	t.showInfo()
}

func (t *tui) getSupervisorLayout() *tview.Flex {
	flex := tview.NewFlex().SetDirection(tview.FlexColumn)
	flex.AddItem(t.getGroupProcessListLayout(), 0, 2, false)
	flex.AddItem(t.getProcessInfoLayout(), 0, 5, false)

	return flex
}

func (t *tui) getGroupProcessListLayout() *tview.Flex {
	groupProcessLists := tview.NewFlex().SetDirection(tview.FlexRow)
	groupProcessLists.AddItem(t.groupLayout, 0, 1, false)
	groupProcessLists.AddItem(t.processLayout, 0, 1, false)

	return groupProcessLists
}

func (t *tui) getProcessInfoLayout() *tview.Flex {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.AddItem(t.getTabsLayout(), 0, 1, false)
	flex.AddItem(t.infoTextView, 0, 12, false)

	return flex
}

func (t *tui) getTabsLayout() *tview.Flex {
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

	return flex
}

func (t *tui) setInfoTextView() {
	textView := tview.NewTextView()
	textView.SetBorder(true)
	textView.SetWrap(true)
	textView.SetFocusFunc(nil)

	t.infoTextView = textView
}

func createButton(name string, handler func()) *tview.Button {
	button := tview.NewButton(name).SetSelectedFunc(handler)
	button.SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)

	return button
}
