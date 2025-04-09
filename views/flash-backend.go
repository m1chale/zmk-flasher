package views

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/new-er/zmk-flasher/platform"
)

type FlashBackend struct {
	isListeningForBlockDevices bool
}

func (m FlashBackend) Init() tea.Cmd {
	return nil
}

func (m FlashBackend) Update(msg tea.Msg) (FlashBackend, tea.Cmd) {
	cmds := []tea.Cmd{}
	switch msg := msg.(type) {
	case BlockDevicesReceivedMsg:
		if m.isListeningForBlockDevices {
			cmds = append(cmds, getBlockDevicesEveryCmd())
		}
	case StartBlockDeviceListenerMsg:
		m.isListeningForBlockDevices = true
	case StopBlockDeviceListener:
		m.isListeningForBlockDevices = false
	}

	return m, tea.Batch(cmds...)
}

func getBlockDevicesEveryCmd() tea.Cmd {
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

type BlockDevicesReceivedMsg struct {
	BlockDevices []platform.BlockDevice
}

type StartBlockDeviceListenerMsg struct{}
type StopBlockDeviceListener struct{}

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
