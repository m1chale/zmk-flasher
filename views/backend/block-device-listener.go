package backend

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/new-er/zmk-flasher/platform"
	"github.com/new-er/zmk-flasher/slices"
	s "slices"
)

type BlockDeviceListener struct {
	prevBlockDevices []platform.BlockDevice
}

func (m BlockDeviceListener) Init() tea.Cmd {
	return getBlockDevicesEveryCmd()
}

func (m BlockDeviceListener) Update(msg tea.Msg) (BlockDeviceListener, tea.Cmd) {
	switch msg := msg.(type) {
	case blockDevicesReceivedMsg:
		if s.EqualFunc(msg.BlockDevices, m.prevBlockDevices, platform.CompareBlockDevice) {
			return m, getBlockDevicesEveryCmd()
		}
		added := []platform.BlockDevice{}
		removed := []platform.BlockDevice{}
		if m.prevBlockDevices != nil {
			added = slices.GetAddedElements(m.prevBlockDevices, msg.BlockDevices, platform.CompareBlockDevice)
			removed = slices.GetRemovedElements(m.prevBlockDevices, msg.BlockDevices, platform.CompareBlockDevice)
		}
		m.prevBlockDevices = msg.BlockDevices
		blockDevicesChangedMsg := BlockDevicesChangedMsg{
			BlockDevices: msg.BlockDevices,
			Added:        added,
			Removed:      removed,
		}
		return m, tea.Batch(getBlockDevicesEveryCmd(), Cmd(blockDevicesChangedMsg))
	}

	return m, nil
}

func getBlockDevicesEveryCmd() tea.Cmd {
	return tea.Every(time.Millisecond * 50, func(t time.Time) tea.Msg {
		devices, err := platform.Operations.GetBlockDevices()
		if err != nil {
			return err
		}

		return blockDevicesReceivedMsg{
			BlockDevices: devices,
		}
	})
}

type blockDevicesReceivedMsg struct {
	BlockDevices []platform.BlockDevice
}

type BlockDevicesChangedMsg struct {
	BlockDevices []platform.BlockDevice
	Added        []platform.BlockDevice
	Removed      []platform.BlockDevice
}
