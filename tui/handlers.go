package tui

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
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

func (t *Tui) tailStdoutLogs() {
	t.infoTextView.Clear()

	info := t.groupToProcessesInfo[t.selectedGroup][t.selectedProcess]
	if info.StdoutLogfile == "" {
		return
	}

	// TODO: configure this
	offset, length := 0, 4000

	go func() {
		for {
			logSegment, err := t.client.TailProcessStdoutLog(t.selectedProcess, offset, length)
			if err != nil {
				log.Errorf("failed to tail stdout logs because %s", err.Error())
				t.app.Stop()
			}

			if logSegment == nil {
				continue
			}

			t.app.QueueUpdateDraw(func() {
				_, err := t.infoTextView.Write([]byte(logSegment.Payload))
				if err != nil {
					log.Errorf("failed to write stdout logs to the screen because %s", err.Error())
					t.app.Stop()
				}
			})

			offset = int(logSegment.Offset)

			if logSegment.Overflow {
				// TODO: configure this
				// nolint:gomnd
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()
}

func (t *Tui) tailStderrLogs() {
	t.infoTextView.Clear()

	info := t.groupToProcessesInfo[t.selectedGroup][t.selectedProcess]
	if info.StdoutLogfile == "" {
		return
	}

	// TODO: configure this
	offset, length := 0, 4000

	go func() {
		for {
			logSegment, err := t.client.TailProcessStderrLog(t.selectedProcess, offset, length)
			if err != nil {
				log.Errorf("failed to tail stderr logs because %s", err.Error())
				t.app.Stop()
			}

			if logSegment == nil {
				continue
			}

			t.app.QueueUpdateDraw(func() {
				_, err := t.infoTextView.Write([]byte(logSegment.Payload))
				if err != nil {
					log.Errorf("failed to write stderr logs to the screen because %s", err.Error())
					t.app.Stop()
				}
			})

			offset = int(logSegment.Offset)

			if logSegment.Overflow {
				// TODO: configure this
				// nolint:gomnd
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()
}
