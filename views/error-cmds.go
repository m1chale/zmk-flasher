package views

import tea "github.com/charmbracelet/bubbletea"

func ErrorCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return err
	}
}
