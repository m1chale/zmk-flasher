package views

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/new-er/zmk-flasher/views/backend"
)

var (
	selectedKeyboardStyle   = lipgloss.NewStyle().Margin(0, 1).Bold(true).Foreground(backend.PrimaryColor)
	unselectedKeyboardStyle = lipgloss.NewStyle().Margin(0, 1).Foreground(lipgloss.Color("240"))
)

type FlashView struct {
	blockDeviceListener backend.BlockDeviceListener

	automationView             AutomationView
	centralKeyboardHalfView    KeyboardHalfView
	peripheralKeyboardHalfView KeyboardHalfView
	selectedKeyboardHalf       backend.KeyboardHalfRole

	dryRun bool
}

func NewFlashView(centralBootloaderFile, peripheralBootloaderFile string, centralMountPoint, peripheralMountPoint *string, dryRun bool) FlashView {
	return FlashView{
		automationView: NewAutomationView(),
		centralKeyboardHalfView: NewKeyboardHalfView(
			backend.Central,
			centralBootloaderFile,
			centralMountPoint,
			dryRun,
		),
		peripheralKeyboardHalfView: NewKeyboardHalfView(
			backend.Peripheral,
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
		f.automationView.Init(),
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

	model, cmd := f.automationView.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	f.automationView = model.(AutomationView)

	model, cmd = f.peripheralKeyboardHalfView.Update(msg)
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
			if f.CanToggleKeyboardHalf() && f.automationView.currentAutomationStrategyIndex == -1 {
				f.selectedKeyboardHalf = f.selectedKeyboardHalf.Toggle()
			}
		case "l":
			if f.CanToggleKeyboardHalf() && f.automationView.currentAutomationStrategyIndex == -1 {
				f.selectedKeyboardHalf = f.selectedKeyboardHalf.Toggle()
			}
		case "0":
			if f.CanToggleKeyboardHalf() && f.automationView.currentAutomationStrategyIndex == -1 {
				cmds = append(cmds, backend.Cmd(StartAutomationMsg{strategyIndex: 0}))
			}
		case "1":
			if f.CanToggleKeyboardHalf() && f.automationView.currentAutomationStrategyIndex == -1 {
				cmds = append(cmds, backend.Cmd(StartAutomationMsg{strategyIndex: 1}))
			}
		case "enter":
			if f.automationView.currentAutomationStrategyIndex == -1 {
				keyboardHalf := f.getSelectedKeyboardHalf()
				cmds = append(cmds, backend.Cmd(NextStepMsg{role: keyboardHalf.role}))
			}
		}
	case StartInteractiveMountMsg:
		f.selectedKeyboardHalf = msg.role
	case StartFlashMsg:
		f.selectedKeyboardHalf = msg.role
	case FlashFinishedMsg:
		if f.centralKeyboardHalfView.step == Done && f.peripheralKeyboardHalfView.step == Done {
			return f, tea.Quit
		}
	case error:
		println("error")
		println(msg.Error())
		return f, tea.Quit
	}

	return f, tea.Batch(cmds...)
}

func (f FlashView) CanToggleKeyboardHalf() bool {
	return f.centralKeyboardHalfView.CanUnselect() && f.peripheralKeyboardHalfView.CanUnselect()
}
func (f FlashView) getSelectedKeyboardHalf() KeyboardHalfView {
	if f.selectedKeyboardHalf == backend.Central {
		return f.centralKeyboardHalfView
	}
	return f.peripheralKeyboardHalfView
}

func (f FlashView) View() string {
	b := strings.Builder{}
	b.WriteString("Zmk Flasher ")
	if f.dryRun {
		b.WriteString("(DryRun)")
	}
	b.WriteString("\n")
	b.WriteString("Press 'q' to quit\n")
	if f.CanToggleKeyboardHalf() && f.automationView.currentAutomationStrategyIndex != 1 {
		b.WriteString("Press 'h'/'l' to switch between keyboard halves\n")
	}
	b.WriteString(f.automationView.View())
	b.WriteString("----------------\n")

	b.WriteString(
		lipgloss.JoinHorizontal(
			lipgloss.Center,
			getKeyboardHalfViewStyle(backend.Central, f.selectedKeyboardHalf).Render(f.centralKeyboardHalfView.View()),
			getKeyboardHalfViewStyle(backend.Peripheral, f.selectedKeyboardHalf).Render(f.peripheralKeyboardHalfView.View()),
		),
	)

	return b.String()
}

func getKeyboardHalfViewStyle(role, selectedKeyboardHalf backend.KeyboardHalfRole) lipgloss.Style {
	if role == selectedKeyboardHalf {
		return selectedKeyboardStyle
	}
	return unselectedKeyboardStyle
}
