# Lazysupervisor
## Lazy way of handling supervisor

A simple terminal UI for supervisor, written in Go with the [tview](https://github.com/rivo/tview) library.
It is inspired by the TUI built for docker and docker-compose called [lazydocker](https://github.com/jesseduffield/lazydocker).

## Features

- View groups and processes
- Start, stop, remove group
- Start, stop process
- View process info, stdout and stderr logs


## TODO
- [ ] Add button to reload superviord
- [ ] Add bulk commands (remove all groups, start all processes, etc)