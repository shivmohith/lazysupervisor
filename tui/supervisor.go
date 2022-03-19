package tui

import "github.com/shivmohith/go-supervisord"

func (t *Tui) getGroupProcessesMap() map[string]map[string]supervisord.ProcessInfo {
	groupToProcesses := make(map[string]map[string]supervisord.ProcessInfo)

	// TODO: do not ignore error
	processes, _ := t.client.GetAllProcessInfo()

	for _, p := range processes {
		if _, ok := groupToProcesses[p.Group]; !ok {
			groupToProcesses[p.Group] = make(map[string]supervisord.ProcessInfo)
		}

		groupToProcesses[p.Group][p.Name] = p
	}

	return groupToProcesses
}
