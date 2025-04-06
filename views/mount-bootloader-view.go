package views

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"zmk-flasher/platform"
	"zmk-flasher/slices"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	mountCmdRunning bool
)

type MountBootloaderView struct {
	bootloaderName               string
	devicesInitialized           bool
	initialDevices               []platform.BlockDevice
	deviceCandidates             []platform.BlockDevice
	selectedDeviceCandidateIndex int
	mountPath                    string
}

func NewMountBootloaderView(bootloaderName string) MountBootloaderView {
	return MountBootloaderView{
		devicesInitialized: false,
		bootloaderName: bootloaderName,
	}
}

func (f MountBootloaderView) EnsureMountPathExists() error {
	if f.mountPath == "" {
		return errors.New("mount path is empty")
	}
	if _, err := os.Stat(f.mountPath); os.IsNotExist(err) {
		return errors.New("mount path does not exist")
	}
	return nil
}

func (f MountBootloaderView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j":
			if f.selectedDeviceCandidateIndex < len(f.deviceCandidates)-1 {
				f.selectedDeviceCandidateIndex++
			}
		case "k":
			if f.selectedDeviceCandidateIndex > 0 {
				f.selectedDeviceCandidateIndex--
			}
		case "enter":
			if len(f.deviceCandidates) > 0 {
				return f, mountBootloaderCmd(f.bootloaderName, f.deviceCandidates[f.selectedDeviceCandidateIndex])
			}
		}
	case BlockDevicesReceivedMsg:
		if !f.devicesInitialized {
			f.initialDevices = msg.Devices
			f.devicesInitialized = true
		}
		addedDevices := slices.GetAddedElements(f.initialDevices, msg.Devices, func(a, b platform.BlockDevice) bool {
			return strings.EqualFold(a.Name, b.Name)
		})

		if len(addedDevices) == 0 {
			f.initialDevices = msg.Devices
			return f, getBlockDevicesCmd()
		}
		if len(addedDevices) == 1 {
			return f, mountBootloaderCmd(f.bootloaderName, addedDevices[0])
		}
		f.deviceCandidates = addedDevices
		return f, nil
	case BootloaderMountedMsg:
		f.mountPath = msg.Device.MountPoints[len(msg.Device.MountPoints)-1]
		f.initialDevices = nil
		return f, NextStepCmd()
	}
	if f.initialDevices == nil {
		return f, getBlockDevicesCmd()
	}
	if f.mountPath != "" {
		return f, NextStepCmd()
	}
	return f, nil
}

func (f MountBootloaderView) View() string {
	b := strings.Builder{}
	if f.initialDevices == nil {
		b.WriteString(f.bootloaderName + " :Enumerating connected devices. Please wait...\n")
	} else if len(f.deviceCandidates) == 0 {
		currentDevices := len(f.initialDevices)
		currentDevicesStr := strconv.FormatInt(int64(currentDevices), 10)
		b.WriteString("Please connect the " + f.bootloaderName + " bootloader (current devices: " + currentDevicesStr + ")\n")
	} else if len(f.deviceCandidates) > 0 {
		b.WriteString("Select the " + f.bootloaderName + " bootloader device\n")
		for i, device := range f.deviceCandidates {
			if i == f.selectedDeviceCandidateIndex {
				b.WriteString(">")
			}
			b.WriteString(device.Name)
			b.WriteString(" - ")
			b.WriteString(device.Label)
			b.WriteString("\n")
		}
	}

	return b.String()
}

func (f MountBootloaderView) Init() tea.Cmd {
	return nil
}

func getBlockDevicesCmd() tea.Cmd {
	return func() tea.Msg {
		devices, err := platform.Operations.GetBlockDevices()
		if err != nil {
			return errors.Join(errors.New("error during block device enumeration"), err)
		}
		return BlockDevicesReceivedMsg{Devices: devices}
	}
}

func mountBootloaderCmd(bootloaderName string, device platform.BlockDevice) tea.Cmd {
	if mountCmdRunning {
		return nil
	}
	mountCmdRunning = true
	return func() tea.Msg {
		device, err := platform.Operations.MountBlockDevice(device)
		if err != nil {
			mountCmdRunning = false
			return errors.Join(errors.New("error during mount"), err)
		}
		mountCmdRunning = false
		return BootloaderMountedMsg{
			bootloaderName: bootloaderName,
			Device:         device,
		}
	}
}

type BlockDevicesReceivedMsg struct {
	Devices []platform.BlockDevice
}

type BootloaderMountedMsg struct {
	bootloaderName string
	Device         platform.BlockDevice
}
