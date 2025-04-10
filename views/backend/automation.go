package backend

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	highlightedText = lipgloss.NewStyle().Foreground(PrimaryColor).Bold(true)
)

type AutomationAction int

const (
	Mount AutomationAction = iota
	Flash
)

func (a AutomationAction) String() string {
	switch a {
	case Mount:
		return "Mount"
	case Flash:
		return "Flash"
	}
	return "Unknown action"
}

var (
	AutomationStrategies = []AutomationStrategy{
		AutomationStrategy{
			Steps: []AutomationStep{
				NewAutomationStep(Mount, Peripheral),
				NewAutomationStep(Flash, Peripheral),
				NewAutomationStep(Mount, Central),
				NewAutomationStep(Flash, Central),
			},
		},
		AutomationStrategy{
			Steps: []AutomationStep{
				NewAutomationStep(Mount, Peripheral),
				NewAutomationStep(Mount, Central),
				NewAutomationStep(Flash, Peripheral),
				NewAutomationStep(Flash, Central),
			},
		},
	}
)

type AutomationStrategy struct {
	Steps []AutomationStep
}


func (s AutomationStrategy) String(selectedStep int) string {
	b := strings.Builder{}
	for i, step := range s.Steps {
		if i == selectedStep {
			b.WriteString(highlightedText.Render(step.String()))
		} else {
			b.WriteString(step.String())
		}
		if i < len(s.Steps)-1 {
			b.WriteString(" -> ")
		}
	}
	return b.String()
}

type AutomationStep struct {
	Action AutomationAction
	Side   KeyboardHalfRole
}

func (s AutomationStep) String() string {
	b := strings.Builder{}
	switch s.Action {
	case Flash:
		b.WriteString("Flash ")
	case Mount:
		b.WriteString("Mount ")
	}
	b.WriteString(s.Side.String())
	return b.String()
}

func NewAutomationStep(action AutomationAction, side KeyboardHalfRole) AutomationStep {
	return AutomationStep{
		Action: action,
		Side:   side,
	}
}
