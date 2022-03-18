package tui

import (
	"log"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	gosupervisord "github.com/shivmohith/go-supervisord"
	"github.com/shivmohith/tui-supervisor/supervisord"
)

type tui struct {
	client supervisord.Client

	groupToProcessesInfo map[string][]gosupervisord.ProcessInfo

	app    *tview.Application
	layout *tview.Flex

	headerLayout *tview.TextView
	footerLayout *tview.TextView

	groups        map[string]struct{}
	groupLayout   *tview.List
	processLayout *tview.List
}

func New() *tui {
	t := new(tui)

	client, err := supervisord.NewClient()
	if err != nil {
		log.Fatalf("failed to get new supervisord client because %v", err)
	}

	t.client = client
	t.app = tview.NewApplication()

	t.setHeader()

	t.groupToProcessesInfo = t.getGroupProcessesMap()
	t.setGroupLayout()
	t.setProcessLayout()
	t.setLayout()

	return t
}

func (t *tui) refreshGroupsAndProcesses() {
	gTop := t.getGroupProcessesMap()
	for g, p := range gTop {
		gCopy := g
		if _, ok := t.groupToProcessesInfo[gCopy]; !ok {
			t.groupToProcessesInfo[gCopy] = p
			t.groupLayout.AddItem(gCopy, "", 0, func() {
				t.setProcesses(gCopy)
			})
		}
	}
}

func (t *tui) Start() {
	runEvery(1*time.Second, func() {
		t.refreshGroupsAndProcesses()
	})

	if err := t.app.SetRoot(t.layout, true).SetFocus(t.layout).EnableMouse(true).Run(); err != nil {
		log.Fatalf("failed to run the tui app because %v", err)
	}
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
			t.setProcesses(gCopy)
		})
	}

	t.groupLayout = list
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

	for _, p := range t.groupToProcessesInfo[group] {
		t.processLayout.AddItem(p.Name, "", 0, nil)
	}
}

func (t *tui) getSupervisorLayout() *tview.Flex {
	groupProcessLists := tview.NewFlex().SetDirection(tview.FlexRow)
	groupProcessLists.AddItem(t.groupLayout, 0, 1, false)
	groupProcessLists.AddItem(t.processLayout, 0, 1, false)

	flex := tview.NewFlex().SetDirection(tview.FlexColumn)
	flex.AddItem(groupProcessLists, 0, 2, false)
	flex.AddItem(tview.NewBox().SetBorder(false).SetTitle("Main info"), 0, 5, false)

	return flex
}
