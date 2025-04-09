package views

import (
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/new-er/zmk-flasher/views/backend"
)

var (
	selectedKeyboardStyle   = lipgloss.NewStyle().Margin(0, 1).Bold(true).Foreground(lipgloss.Color("205"))
	unselectedKeyboardStyle = lipgloss.NewStyle().Margin(0, 1).Foreground(lipgloss.Color("240"))
)

type FlashView struct {
	blockDeviceListener backend.BlockDeviceListener

	centralKeyboardHalfView    KeyboardHalfView
	peripheralKeyboardHalfView KeyboardHalfView
	selectedKeyboardHalf       KeyboardHalfRole

	dryRun bool
}

func NewFlashView(centralBootloaderFile, peripheralBootloaderFile string, centralMountPoint, peripheralMountPoint *string, dryRun bool) FlashView {
	return FlashView{
		centralKeyboardHalfView: NewKeyboardHalfView(
			Central,
			centralBootloaderFile,
			centralMountPoint,
			dryRun,
		),
		peripheralKeyboardHalfView: NewKeyboardHalfView(
			Peripheral,
			peripheralBootloaderFile,
			peripheralMountPoint,
			dryRun,
		),
		dryRun: dryRun,
	}
}

func (f FlashView) Init() tea.Cmd {
	return tea.Batch(
		f.blockDeviceListener.Init(),
		f.centralKeyboardHalfView.Init(),
		f.peripheralKeyboardHalfView.Init(),
	)
}

func (f FlashView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	backendModel, backendCmd := f.blockDeviceListener.Update(msg)
	if backendCmd != nil {
		cmds = append(cmds, backendCmd)
	}
	f.blockDeviceListener = backendModel

	model, cmd := f.peripheralKeyboardHalfView.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	f.peripheralKeyboardHalfView = model.(KeyboardHalfView)

	model, cmd = f.centralKeyboardHalfView.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	f.centralKeyboardHalfView = model.(KeyboardHalfView)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return f, tea.Quit
		case "h":
			if f.CanToggleKeyboardHalf() {
				f.selectedKeyboardHalf = f.selectedKeyboardHalf.Toggle()
			}
		case "l":
			if f.CanToggleKeyboardHalf() {
				f.selectedKeyboardHalf = f.selectedKeyboardHalf.Toggle()
			}
		case "enter":
			if f.selectedKeyboardHalf == Central {
				m, cmd := f.centralKeyboardHalfView.NextStep()
				f.centralKeyboardHalfView = m
				cmds = append(cmds, cmd)
			} else {
				m, cmd := f.peripheralKeyboardHalfView.NextStep()
				f.peripheralKeyboardHalfView = m
				cmds = append(cmds, cmd)
			}
		}
	case NextStepMsg:
		if f.selectedKeyboardHalf == Central {
			m, cmd := f.centralKeyboardHalfView.NextStep()
			f.centralKeyboardHalfView = m
			cmds = append(cmds, cmd)
		} else {
			m, cmd := f.peripheralKeyboardHalfView.NextStep()
			f.peripheralKeyboardHalfView = m
			cmds = append(cmds, cmd)
		}
	case error:
		println("error")
		println(msg.Error())
		return f, tea.Quit
	}

	(&f).updateKeyboardHalfIsSelected()

	return f, tea.Batch(cmds...)
}

func (f *FlashView) updateKeyboardHalfIsSelected() {
	if f.selectedKeyboardHalf == Central {
		f.centralKeyboardHalfView = f.centralKeyboardHalfView.SetIsSelected(true)
		f.peripheralKeyboardHalfView = f.peripheralKeyboardHalfView.SetIsSelected(false)
	} else {
		f.centralKeyboardHalfView = f.centralKeyboardHalfView.SetIsSelected(false)
		f.peripheralKeyboardHalfView = f.peripheralKeyboardHalfView.SetIsSelected(true)
	}
}

func (f FlashView) CanToggleKeyboardHalf() bool {
	return f.centralKeyboardHalfView.CanUnselect() && f.peripheralKeyboardHalfView.CanUnselect()
}

func (f FlashView) View() string {
	b := strings.Builder{}
	b.WriteString("Zmk Flasher\n")
	b.WriteString("Dry Run: " + strconv.FormatBool(f.dryRun) + "\n")
	b.WriteString("Press 'q' to quit\n")
	if f.CanToggleKeyboardHalf() {
		b.WriteString("Press 'h'/'l' to switch between keyboard halves\n")
	}
	b.WriteString("----------------\n")

	b.WriteString(
		lipgloss.JoinHorizontal(
			lipgloss.Center,
			getKeyboardHalfViewStyle(Central, f.selectedKeyboardHalf).Render(f.centralKeyboardHalfView.View()),
			getKeyboardHalfViewStyle(Peripheral, f.selectedKeyboardHalf).Render(f.peripheralKeyboardHalfView.View()),
		),
	)

	return b.String()
}

func getKeyboardHalfViewStyle(role, selectedKeyboardHalf KeyboardHalfRole) lipgloss.Style {
	if role == selectedKeyboardHalf {
		return selectedKeyboardStyle
	}
	return unselectedKeyboardStyle
}

func NextStepCmd() tea.Cmd {
	return func() tea.Msg {
		return NextStepMsg{}
	}
}

type NextStepMsg struct{}
