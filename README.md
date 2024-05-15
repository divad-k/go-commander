# Go Manager
## Overview
The Go manager application is a comprehensive file management tool developed in Go. Users can navigate through directories, copy, move, and delete files. Built with [Tview](https://github.com/rivo/tview), [Tcell](https://github.com/gdamore/tcell). Currently undergoing testing.

<p align="center">
  <img src="https://github.com/divad-k/go-commander/blob/main/screenshot.png" width="350" alt="screenshot">
</p>

## Instructions

Clone the repository:
```bash
git clone https://github.com/divad-k/go-manager.git
```
Navigate to the cloned directory:
```bash
cd go-manager
```
Build the project using Go compiler:
```bash
go build .
```
## Usage
- Use arrow keys to navigate through the lists.
- Press Enter to select an item or directory.
- Press Tab to switch focus between lists.
- Use keys c for copy, m for move, d for delete, i for info, and q to quit.
