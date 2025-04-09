package views

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/new-er/zmk-flasher/platform"
)

type FlashBackend struct {
	shouldUpdateBlockDevices bool
}

func (m FlashBackend) Init() tea.Cmd {
	return nil
}

func (m FlashBackend) Update(msg tea.Msg) (FlashBackend, tea.Cmd) {
	switch msg := msg.(type) {
	case BlockDevicesReceivedMsg:
		if m.shouldUpdateBlockDevices {
			return m, updateBlockDevicesEveryCmd()
		}
	case changeUpdateBlockDevicesMsg:
		m.shouldUpdateBlockDevices = msg.ShouldUpdateBlockDevices
		return m, updateBlockDevicesEveryCmd()
	}
	return m, nil
}

func updateBlockDevicesEveryCmd() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		devices, err := platform.Operations.GetBlockDevices()
		if err != nil {
			return err
		}
		return BlockDevicesReceivedMsg{
			BlockDevices: devices,
		}
	})
}

func MountBlockDeviceCmd(device platform.BlockDevice) tea.Cmd {
	return func() tea.Msg {
		mountedDevice, err := platform.Operations.MountBlockDevice(device)
		if err != nil {
			return err
		}

		return BlockDeviceMountedMsg{
			BlockDevice: mountedDevice,
		}
	}
}

type BlockDeviceMountedMsg struct {
	BlockDevice platform.BlockDevice
}

func ChangeUpdateBlockDevicesCmd(shouldUpdate bool) tea.Cmd {
	return func() tea.Msg {
		return changeUpdateBlockDevicesMsg{
			ShouldUpdateBlockDevices: shouldUpdate,
		}
	}
}

type changeUpdateBlockDevicesMsg struct {
	ShouldUpdateBlockDevices bool
}

type BlockDevicesReceivedMsg struct {
	BlockDevices []platform.BlockDevice
}
