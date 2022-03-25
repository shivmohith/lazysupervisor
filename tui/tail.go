package tui

import (
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

	go func() {
		for t.tabsLayout.GetItem(1).HasFocus() {
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
		for t.tabsLayout.GetItem(2).HasFocus() {
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
