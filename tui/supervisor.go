package tui

import "github.com/shivmohith/go-supervisord"

func (t *tui) getGroupProcessesMap() map[string][]supervisord.ProcessInfo {
	groupToProcessesInfo := make(map[string][]supervisord.ProcessInfo)

	// TODO: do not ignore error
	processes, _ := t.client.GetAllProcessInfo()

	for _, p := range processes {
		groupToProcessesInfo[p.Group] = append(groupToProcessesInfo[p.Group], p)
	}

	return groupToProcessesInfo
}
