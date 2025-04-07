package views

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/new-er/zmk-flasher/files"
)

type KeyboardHalfRole int

const (
	Central KeyboardHalfRole = iota
	Peripheral
)

func (k KeyboardHalfRole) String() string {
	switch k {
	case Central:
		return "Central"
	case Peripheral:
		return "Peripheral"
	default:
		return "Unknown"
	}
}

type KeyboardHalfView struct {
	role           KeyboardHalfRole
	bootloaderFile string
	mountPath      *string
}

func NewKeyboardHalfView(role KeyboardHalfRole, bootloaderFile string, mountPath *string) KeyboardHalfView {
	return KeyboardHalfView{
		role:           role,
		bootloaderFile: bootloaderFile,
		mountPath:      mountPath,
	}
}

func (k KeyboardHalfView) Init() tea.Cmd {
	return nil
}

func (k KeyboardHalfView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return k, nil
}

func (k KeyboardHalfView) View() string {
	b := strings.Builder{}
	b.WriteString(k.role.String())
	b.WriteString("\n")
	b.WriteString("🗎 :")
	b.WriteString(files.EllipsisFront(k.bootloaderFile, 40))
	b.WriteString("\n")

	b.WriteString("󱊞 : ")
	if k.mountPath != nil {
		b.WriteString(*k.mountPath)
	} else {
		b.WriteString("Not mounted")
	}
	b.WriteString("\n")

	return b.String()
}

func ellipsisFront(s string, maxLen int) string {
		runes := []rune(s)
		if len(runes) <= maxLen {
				return s
		}
		if maxLen < 3 {
				maxLen = 3
		}
		return "..." + string(runes[len(runes)-maxLen+3:])
}
