package views

import (
	"errors"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/new-er/zmk-flasher/files"
	"github.com/new-er/zmk-flasher/views/backend"
)

type KeyboardHalfRole int

const (
	Central KeyboardHalfRole = iota
	Peripheral
)

func (k KeyboardHalfRole) Toggle() KeyboardHalfRole {
	switch k {
	case Central:
		return Peripheral
	case Peripheral:
		return Central
	default:
		return k
	}
}

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

type FlashStep int

const (
	Unmounted FlashStep = iota
	SelectDevice
	Mounting
	ReadyToFlash
	Flashing
	Done
)

type KeyboardHalfView struct {
	step           FlashStep
	role           KeyboardHalfRole
	bootloaderFile string
	mountPath      *string
	isSelected     bool

	dryRun bool
}

func NewKeyboardHalfView(role KeyboardHalfRole, bootloaderFile string, mountPath *string, dryRun bool) KeyboardHalfView {
	return KeyboardHalfView{
		role:           role,
		bootloaderFile: bootloaderFile,
		mountPath:      mountPath,
	}
}

func (k KeyboardHalfView) SetIsSelected(isSelected bool) KeyboardHalfView {
	k.isSelected = isSelected
	return k
}

func (k KeyboardHalfView) CanUnselect() bool {
	return k.step == Unmounted || k.step == ReadyToFlash || k.step == Done
}

func (k KeyboardHalfView) NextStep() (KeyboardHalfView, tea.Cmd) {
	k.step++
	if k.step > Done {
		k.step = Done
	}
	if k.step == Flashing {
		return k, backend.CopyFileCmd(k.bootloaderFile, (*k.mountPath)+"/firmware.uf2", k.dryRun)
	}

	return k, nil
}

func (k KeyboardHalfView) Init() tea.Cmd {
	return nil
}

func (k KeyboardHalfView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case backend.BlockDevicesChangedMsg:
		if k.step == SelectDevice {
			if len(msg.Added) == 0 {
				return k, nil
			}
			if len(msg.Added) > 1 {
				return k, backend.Cmd(errors.New("multiple devices added"))
			}

			return k, tea.Batch(
				backend.MountBlockDeviceCmd(msg.Added[0]),
				NextStepCmd(),
			)
		}
	case backend.BlockDeviceMountedMsg:
		if k.step == Mounting {
			k.step = ReadyToFlash
			k.mountPath = &msg.BlockDevice.MountPoints[0]
			return k, nil
		}
	case backend.FileCopiedMsg:
		if k.step == Flashing {
			k.step = Done
			return k, nil
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
		b.WriteString("")
	}
	b.WriteString("\n")

	if k.isSelected {
		switch k.step {
		case Unmounted:
			b.WriteString("Press [Enter] to mount bootloader.\n")
		case SelectDevice:
			b.WriteString("Please connect the ")
			b.WriteString(k.role.String())
			b.WriteString(" controller.\n")
		case Mounting:
			b.WriteString("Mounting, please wait...\n")
		case ReadyToFlash:
			b.WriteString("Press [Enter] to flash the bootloader.\n")
		case Flashing:
			b.WriteString("Flashing, please wait...\n")
		case Done:
			b.WriteString("\n")
		}
	}
	return b.String()
}
