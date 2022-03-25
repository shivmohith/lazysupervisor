package tui

import (
	"fmt"
	"time"
)

const infoTemplate = `Name: %v
Group: %v
Pid: %v
Start time: %v
Stop time: %v
State: %v
State name: %v
Spawn error: %v
Exit status: %v
Stdout log file: %v
Stderr log file: %v
`

func (t *Tui) showInfo() {
	pInfo := t.groupToProcessesInfo[t.selectedGroup][t.selectedProcess]

	info := fmt.Sprintf(
		infoTemplate,
		pInfo.Name,
		pInfo.Group,
		pInfo.Pid,
		time.Unix(int64(pInfo.Start), 0),
		time.Unix(int64(pInfo.Stop), 0),
		pInfo.State,
		pInfo.StateName,
		pInfo.SpawnErr,
		pInfo.ExitStatus,
		pInfo.StdoutLogfile,
		pInfo.StderrLogfile,
	)

	t.infoTextView.SetText(info)
}
