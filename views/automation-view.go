package views

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/new-er/zmk-flasher/views/backend"
)

type AutomationView struct {
	currentAutomationStrategyIndex int
	currentAutomationStep          int
}

func NewAutomationView() AutomationView {
	return AutomationView{
		currentAutomationStrategyIndex: -1,
		currentAutomationStep:          -1,
	}
}

func (a AutomationView) Init() tea.Cmd {
	return nil
}

func (a AutomationView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case StartAutomationMsg:
		if a.currentAutomationStrategyIndex == -1 {
			a.currentAutomationStrategyIndex = msg.strategyIndex
			a.currentAutomationStep = 0
			return a, backend.Cmd(updateAutomationMsg{})
		}
	case updateAutomationMsg:
		strategy := backend.AutomationStrategies[a.currentAutomationStrategyIndex]
		if a.currentAutomationStep >= len(strategy.Steps) {
			a.currentAutomationStep = -1
			a.currentAutomationStrategyIndex = -1
			return a, backend.Cmd(AutomationFinishedMsg{})
		}
		step := strategy.Steps[a.currentAutomationStep]

		switch step.Action {
		case backend.Mount:
			return a, backend.Cmd(StartInteractiveMountMsg{role: step.Side})
		case backend.Flash:
			return a, backend.Cmd(StartFlashMsg{role: step.Side})
		}
	case InteractiveMountFinishedMsg:
		if a.currentAutomationStep == -1 {
			return a, nil
		}
		a.currentAutomationStep++
		return a, backend.Cmd(updateAutomationMsg{})
	case FlashFinishedMsg:
		if a.currentAutomationStep == -1 {
			return a, nil
		}
		a.currentAutomationStep++
		return a, backend.Cmd(updateAutomationMsg{})
	}
	return a, nil
}

func (a AutomationView) View() string {
	if !a.IsAutomationRunning() {
		return ""
	}
	b := strings.Builder{}

	strategy := backend.AutomationStrategies[a.currentAutomationStrategyIndex]
	b.WriteString(strategy.String(a.currentAutomationStep))
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
