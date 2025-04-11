package views

import (
	"errors"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/new-er/zmk-flasher/files"
	"github.com/new-er/zmk-flasher/views/backend"
)

type Step int

const (
	Unmounted Step = iota
	WaitForMount
	Mounting
	ReadyToFlash
	Flashing
	Done
)

type KeyboardHalfView struct {
	step           Step
	role           backend.KeyboardHalfRole
	bootloaderFile string
	mountPath      *string

	dryRun bool

	foundDevices int

	spinner spinner.Model
}

func NewKeyboardHalfView(role backend.KeyboardHalfRole, bootloaderFile string, mountPath *string, dryRun bool) KeyboardHalfView {
	step := Unmounted
	if mountPath != nil {
		step = ReadyToFlash
	}
	s := spinner.New()
	s.Spinner = spinner.Dot
	return KeyboardHalfView{
		role:           role,
		bootloaderFile: bootloaderFile,
		mountPath:      mountPath,
		step:           step,
		spinner:        s,
		dryRun:         dryRun,
	}
}

func (k KeyboardHalfView) CanUnselect() bool {
	return k.step == Unmounted || k.step == ReadyToFlash || k.step == Done
}

func (m KeyboardHalfView) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m KeyboardHalfView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		s, cmd := m.spinner.Update(msg)
		m.spinner = s
		return m, cmd
	case StartInteractiveMountMsg:
		if msg.role != m.role {
			return m, nil
		}
		m.step = WaitForMount
	case StartFlashMsg:
		if msg.role != m.role {
			return m, nil
		}
		m.step = Flashing
		return m, backend.CopyFileCmd(m.bootloaderFile, (*m.mountPath)+"/firmware.uf2", m.dryRun)
	case NextStepMsg:
		if msg.role != m.role {
			return m, nil
		}
		if m.step == Unmounted {
			return m, backend.Cmd(StartInteractiveMountMsg{role: m.role})
		}
		if m.step == ReadyToFlash {
			return m, backend.Cmd(StartFlashMsg{role: m.role})
		}
	case backend.BlockDevicesChangedMsg:
		if m.step == WaitForMount {
			m.foundDevices = len(msg.BlockDevices)
			if len(msg.Added) == 0 {
				return m, nil
			}
			if len(msg.Added) > 1 {
				return m, backend.Cmd(errors.New("multiple devices added"))
			}

			m.step = Mounting
			return m, tea.Batch(
				backend.MountBlockDeviceCmd(msg.Added[0]),
			)
		} else if m.mountPath != nil {
			for _, rem := range msg.Removed {
				for _, remPath := range rem.MountPoints {
					if remPath == *m.mountPath {
						m.mountPath = nil
						if m.step != Done && m.step != Flashing {
							m.step = Unmounted
						}
					}
				}
			}
		}
	case backend.BlockDeviceMountedMsg:
		if m.step == Mounting {
			m.step = ReadyToFlash
			m.mountPath = &msg.BlockDevice.MountPoints[0]
			return m, backend.Cmd(InteractiveMountFinishedMsg{})
		}
	case backend.FileCopiedMsg:
		if m.step == Flashing {
			m.step = Done
			return m, backend.Cmd(FlashFinishedMsg{role: m.role})
		}
	}

	return m, nil
}

func (m KeyboardHalfView) View() string {
	b := strings.Builder{}
	b.WriteString(m.role.String())
	b.WriteString("\n")
	b.WriteString("🗎 :")
	b.WriteString(files.EllipsisFront(m.bootloaderFile, 40))
	b.WriteString("\n")

	b.WriteString("🗁 : ")
	if m.mountPath != nil {
		b.WriteString(files.EllipsisFront(*m.mountPath, 40))
	} else {
		if m.step == Done {
			b.WriteString("")
		}
	}
	b.WriteString("\n")

	switch m.step {
	case Unmounted:
		b.WriteString("Press [Enter] to mount bootloader.\n")
	case WaitForMount:
		b.WriteString(m.spinner.View())
		b.WriteString("Connect the ")
		b.WriteString(m.role.String())
		b.WriteString(" controller. (devices ")
		b.WriteString(strconv.FormatInt(int64(m.foundDevices), 10))
		b.WriteString(")")
		b.WriteString("\n")
	case Mounting:
		b.WriteString(m.spinner.View())
		b.WriteString("Mounting, please wait...\n")
	case ReadyToFlash:
		b.WriteString("Press [Enter] to flash the bootloader.\n")
	case Flashing:
		b.WriteString(m.spinner.View())
		b.WriteString("Flashing, please wait...\n")
	case Done:
		if m.mountPath != nil {
			b.WriteString("\n")
		}
	}
	return b.String()
}

type StartInteractiveMountMsg struct {
	role backend.KeyboardHalfRole
}

type InteractiveMountFinishedMsg struct {
	role backend.KeyboardHalfRole
}

type StartFlashMsg struct {
	role backend.KeyboardHalfRole
}

type FlashFinishedMsg struct {
	role backend.KeyboardHalfRole
}

type NextStepMsg struct {
	role backend.KeyboardHalfRole
}

type KeyboardHalfUnmountedMsg struct {
	role backend.KeyboardHalfRole
}
