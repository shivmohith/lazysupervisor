package tui

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	gosupervisord "github.com/shivmohith/go-supervisord"
	"github.com/shivmohith/tui-supervisor/supervisord"
	log "github.com/sirupsen/logrus"
)

type Tui struct {
	client supervisord.Client

	groupToProcessesInfo map[string]map[string]gosupervisord.ProcessInfo

	app    *tview.Application
	layout *tview.Flex

	headerLayout *tview.TextView
	// footerLayout *tview.TextView

	groupLayout   *tview.List
	processLayout *tview.List
	tabsLayout    *tview.Flex
	infoTextView  *tview.TextView

	selectedGroup   string
	selectedProcess string
}

func New(client supervisord.Client) *Tui {
	t := new(Tui)

	t.client = client

	t.app = tview.NewApplication()

	t.setHeader()

	t.groupToProcessesInfo = t.getGroupProcessesMap()

	t.setGroupLayout()
	t.setProcessLayout()
	t.setTabsLayout()
	t.setInfoTextView()
	t.setLayout()

	t.selectedGroup, _ = t.groupLayout.GetItemText(t.groupLayout.GetCurrentItem())
	t.selectedProcess, _ = t.processLayout.GetItemText(t.processLayout.GetCurrentItem())

	return t
}

func (t *Tui) refreshGroupsAndProcesses() {
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
	for g := range t.groupToProcessesInfo {
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

func (t *Tui) Start() error {
	log.Info("Starting the tui app")

	runEvery(1*time.Second, func() {
		t.refreshGroupsAndProcesses()
	})

	return t.app.SetRoot(t.layout, true).SetFocus(t.layout).EnableMouse(true).Run()
}

func (t *Tui) setHeader() {
	header := tview.NewTextView()
	header.SetText("Welcome to Supervisor TUI")
	header.SetTextAlign(tview.AlignCenter)
	header.SetBorderPadding(1, 1, 1, 1)
	header.SetBackgroundColor(tcell.ColorBlueViolet)

	t.headerLayout = header
}

func (t *Tui) setLayout() {
	mainLayoutProportion := 12

	t.layout = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(t.headerLayout, 0, 1, false).
		AddItem(t.getMainLayout(), 0, mainLayoutProportion, false)
}

func (t *Tui) getMainLayout() *tview.Flex {
	flex := tview.NewFlex().SetDirection(tview.FlexColumn)

	groupProcessListLayoutProportion := 2
	flex.AddItem(t.getGroupProcessListLayout(), 0, groupProcessListLayoutProportion, false)

	processInfoLayoutProportion := 5
	flex.AddItem(t.getProcessInfoLayout(), 0, processInfoLayoutProportion, false)

	return flex
}

func (t *Tui) getGroupProcessListLayout() *tview.Flex {
	groupProcessLists := tview.NewFlex().SetDirection(tview.FlexRow)
	groupProcessLists.AddItem(t.groupLayout, 0, 1, false)
	groupProcessLists.AddItem(t.processLayout, 0, 1, false)

	return groupProcessLists
}

func (t *Tui) getProcessInfoLayout() *tview.Flex {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.AddItem(t.tabsLayout, 0, 1, false)

	infoTextViewProportion := 12
	flex.AddItem(t.infoTextView, 0, infoTextViewProportion, false)

	return flex
}

func (t *Tui) setInfoTextView() {
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
