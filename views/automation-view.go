package views

import (
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/new-er/zmk-flasher/views/backend"
)

type AutomationView struct {
	currentAutomationStrategyIndex int
	currentAutomationStep          int

	spinner          spinner.Model
	isSpinnerVisible bool
}

func NewAutomationView() AutomationView {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return AutomationView{
		currentAutomationStrategyIndex: -1,
		currentAutomationStep:          -1,
		spinner:                        s,
	}
}

func (m AutomationView) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m AutomationView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		spinner, cmd := m.spinner.Update(msg)
		m.spinner = spinner
		return m, cmd
	case StartAutomationMsg:
		if m.currentAutomationStrategyIndex == -1 {
			m.currentAutomationStrategyIndex = msg.strategyIndex
			m.currentAutomationStep = 0
			return m, backend.Cmd(updateAutomationMsg{})
		}
	case updateAutomationMsg:
		strategy := backend.AutomationStrategies[m.currentAutomationStrategyIndex]
		if m.currentAutomationStep >= len(strategy.Steps) {
			m.currentAutomationStep = -1
			m.currentAutomationStrategyIndex = -1
			return m, backend.Cmd(AutomationFinishedMsg{})
		}
		step := strategy.Steps[m.currentAutomationStep]

		switch step.Action {
		case backend.Mount:
			m.isSpinnerVisible = false
			return m, backend.Cmd(StartInteractiveMountMsg{role: step.Side})
		case backend.Flash:
			return m, backend.Cmd(StartFlashMsg{role: step.Side})
		}
	case InteractiveMountFinishedMsg:
		if m.currentAutomationStep == -1 {
			return m, nil
		}
		m.currentAutomationStep++
		m.isSpinnerVisible = true
		return m, backend.Cmd(updateAutomationMsg{})
	case FlashFinishedMsg:
		if m.currentAutomationStep == -1 {
			return m, nil
		}
		m.currentAutomationStep++
		strategy := backend.AutomationStrategies[m.currentAutomationStrategyIndex]
		if m.currentAutomationStep >= len(strategy.Steps) {
			m.isSpinnerVisible = false
		}
		return m, backend.Cmd(updateAutomationMsg{})
	}
	return m, nil
}

func (m AutomationView) View() string {
	if !m.IsAutomationRunning() {
		return ""
	}
	b := strings.Builder{}

	strategy := backend.AutomationStrategies[m.currentAutomationStrategyIndex]
	if m.isSpinnerVisible {
		b.WriteString(m.spinner.View())
	}
	b.WriteString(strategy.String(m.currentAutomationStep))
	b.WriteString("\n")

	return b.String()
}

func (a AutomationView) IsAutomationRunning() bool {
	return a.currentAutomationStrategyIndex != -1
}

type StartAutomationMsg struct {
	strategyIndex int
}
type AutomationFinishedMsg struct{}
type updateAutomationMsg struct{}
