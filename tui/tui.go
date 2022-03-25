package tui

import (
	"time"

	"github.com/rivo/tview"
	gosupervisord "github.com/shivmohith/go-supervisord"
	"github.com/shivmohith/tui-supervisor/supervisord"
)

type Tui struct {
	client supervisord.Client

	groupToProcessesInfo map[string]map[string]gosupervisord.ProcessInfo

	app    *tview.Application
	layout *tview.Flex

	headerLayout *tview.TextView
	footerLayout *tview.TextView

	groupLayout   *tview.List
	processLayout *tview.List
	tabsLayout    *tview.Flex
	infoTextView  *tview.TextView

	selectedGroup   string
	selectedProcess string

	panels                      map[uint32]tview.Primitive
	currentPanelInFocusPosition uint32
}

func New(client supervisord.Client) *Tui {
	t := new(Tui)

	t.client = client
	t.app = tview.NewApplication()

	return t
}

func (t *Tui) BuildLayout() {
	t.setHeader()
	t.setFooter()

	t.groupToProcessesInfo = t.getGroupProcessesMap()

	t.setGroupLayout()
	t.setProcessLayout()
	t.setTabsLayout()
	t.setInfoTextView()

	t.setAppLayout()

	t.selectedGroup, _ = t.groupLayout.GetItemText(t.groupLayout.GetCurrentItem())
	t.selectedProcess, _ = t.processLayout.GetItemText(t.processLayout.GetCurrentItem())

	t.setPanels()
}

func (t *Tui) Start() error {
	t.captureKeyboardEvents()

	runEvery(1*time.Second, func() {
		t.refreshGroupsAndProcesses()
	})

	return t.app.SetRoot(t.layout, true).SetFocus(t.layout).EnableMouse(true).Run()
}

func (t *Tui) refreshGroupsAndProcesses() {
	gTop := t.getGroupProcessesMap()

	t.app.QueueUpdateDraw(func() {
		// Adds new groups created to the group list
		for g, pMap := range gTop {
			gCopy := g

			if _, ok := t.groupToProcessesInfo[gCopy]; !ok {
				t.groupLayout.AddItem(gCopy, "", 0, func() {
					t.handleGroupSelect(gCopy)
				})
			}

			t.groupToProcessesInfo[gCopy] = pMap
		}

		// Remove groups from the group list that have been removed from supervisord
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
	})
}

func (t *Tui) setAppLayout() {
	mainLayoutProportion := 12

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(t.headerLayout, 0, 1, false).
		AddItem(t.getMainLayout(), 0, mainLayoutProportion, false).
		AddItem(t.footerLayout, 0, 1, false)

	t.layout = flex
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

func (t *Tui) setPanels() {
	t.panels = map[uint32]tview.Primitive{
		1: t.groupLayout,
		2: t.processLayout,
		3: t.tabsLayout,
		4: t.infoTextView,
	}

	t.currentPanelInFocusPosition = 1
}
