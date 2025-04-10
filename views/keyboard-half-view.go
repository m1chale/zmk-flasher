package views

import (
	"errors"
	"strconv"
	"strings"

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
}

func NewKeyboardHalfView(role backend.KeyboardHalfRole, bootloaderFile string, mountPath *string, dryRun bool) KeyboardHalfView {
	step := Unmounted
	if mountPath != nil {
		step = ReadyToFlash
	}
	return KeyboardHalfView{
		role:           role,
		bootloaderFile: bootloaderFile,
		mountPath:      mountPath,
		step:           step,
	}
}

func (k KeyboardHalfView) CanUnselect() bool {
	return k.step == Unmounted || k.step == ReadyToFlash || k.step == Done
}

func (k KeyboardHalfView) Init() tea.Cmd {
	return nil
}

func (k KeyboardHalfView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case StartInteractiveMountMsg:
		if msg.role != k.role {
			return k, nil
		}
		k.step = WaitForMount
	case StartFlashMsg:
		if msg.role != k.role {
			return k, nil
		}
		k.step = Flashing
		return k, backend.CopyFileCmd(k.bootloaderFile, (*k.mountPath)+"/firmware.uf2", k.dryRun)
	case NextStepMsg:
		if msg.role != k.role {
			return k, nil
		}
		if k.step == Unmounted {
			return k, backend.Cmd(StartInteractiveMountMsg{role: k.role})
		}
		if k.step == ReadyToFlash {
			return k, backend.Cmd(StartFlashMsg{role: k.role})
		}
	case backend.BlockDevicesChangedMsg:
		if k.step == WaitForMount {
			k.foundDevices = len(msg.BlockDevices)
			if len(msg.Added) == 0 {
				return k, nil
			}
			if len(msg.Added) > 1 {
				return k, backend.Cmd(errors.New("multiple devices added"))
			}

			k.step = Mounting
			return k, tea.Batch(
				backend.MountBlockDeviceCmd(msg.Added[0]),
			)
		} else if k.mountPath != nil {
			for _, rem := range msg.Removed {
				for _, remPath := range rem.MountPoints {
					if remPath == *k.mountPath {
						k.mountPath = nil
						if k.step != Done && k.step != Flashing {
							k.step = Unmounted
						}
					}
				}
			}
		}
	case backend.BlockDeviceMountedMsg:
		if k.step == Mounting {
			k.step = ReadyToFlash
			k.mountPath = &msg.BlockDevice.MountPoints[0]
			return k, backend.Cmd(InteractiveMountFinishedMsg{})
		}
	case backend.FileCopiedMsg:
		if k.step == Flashing {
			k.step = Done
			return k, backend.Cmd(FlashFinishedMsg{role: k.role})
		}
	}

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
		b.WriteString(files.EllipsisFront(*k.mountPath, 40))
	} else {
		if k.step == Done {
			b.WriteString("")
		} else {
			b.WriteString("")
		}
	}
	b.WriteString("\n")

	switch k.step {
	case Unmounted:
		b.WriteString("Press [Enter] to mount bootloader.\n")
	case WaitForMount:
		b.WriteString("Please connect the ")
		b.WriteString(k.role.String())
		b.WriteString(" controller. (current devices ")
		b.WriteString(strconv.FormatInt(int64(k.foundDevices), 10))
		b.WriteString(")")
		b.WriteString("\n")
	case Mounting:
		b.WriteString("Mounting, please wait...\n")
	case ReadyToFlash:
		b.WriteString("Press [Enter] to flash the bootloader.\n")
	case Flashing:
		b.WriteString("Flashing, please wait...\n")
	case Done:
		if k.mountPath != nil {
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
