package tui

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

func (t *Tui) tailStdoutLogs() {
	t.infoTextView.Clear()

	info := t.groupToProcessesInfo[t.selectedGroup][t.selectedProcess]
	if info.StdoutLogfile == "" {
		return
	}

	// TODO: configure this
	offset, length := 0, 4000

	processName := fmt.Sprintf("%s:%s", t.selectedGroup, t.selectedProcess)

	go func() {
		for t.tabsLayout.GetItem(stdoutTabPosition).HasFocus() {
			logSegment, err := t.client.TailProcessStdoutLog(processName, offset, length)
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

	processName := fmt.Sprintf("%s:%s", t.selectedGroup, t.selectedProcess)

	go func() {
		for t.tabsLayout.GetItem(stderrTabPosition).HasFocus() {
			logSegment, err := t.client.TailProcessStderrLog(processName, offset, length)
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
