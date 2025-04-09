package views

import tea "github.com/charmbracelet/bubbletea"

func Cmd(msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}
