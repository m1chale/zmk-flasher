package views

import (
	"errors"
	"strconv"
	"strings"

	"github.com/new-er/zmk-flasher/files"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type flashStep int

const (
	Init flashStep = iota
	MountRightBootloader
	MountLeftBootloader
	FlashRightBootloader
	FlashLeftBootloader
	Done
)

type FlashView struct {
	peripheralKeyboardHalfView KeyboardHalfView
	centralKeyboardHalfView    KeyboardHalfView

	dryRun bool
}

func NewFlashView(centralBootloaderFile, peripheralBootloaderFile string, centralMountPoint, peripheralMountPoint *string, dryRun bool) FlashView {
	return FlashView{
		centralKeyboardHalfView: NewKeyboardHalfView(
			Central,
			centralBootloaderFile,
			centralMountPoint,
		),
		peripheralKeyboardHalfView: NewKeyboardHalfView(
			Peripheral,
			peripheralBootloaderFile,
			peripheralMountPoint,
		),
		dryRun: dryRun,
	}
}

func (f FlashView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	model, cmd := f.peripheralKeyboardHalfView.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	f.peripheralKeyboardHalfView = model.(KeyboardHalfView)

	model, cmd = f.centralKeyboardHalfView.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	f.centralKeyboardHalfView = model.(KeyboardHalfView)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" {
			return f, tea.Quit
		}
	case error:
		println("error")
		println(msg.Error())
		return f, tea.Quit
	}

	return f, tea.Batch(cmds...)
}

func (f FlashView) View() string {
	b := strings.Builder{}
	b.WriteString("Zmk Flasher\n")
	b.WriteString("Dry Run: " + strconv.FormatBool(f.dryRun) + "\n")
	b.WriteString("Press 'q' to quit\n")
	b.WriteString("----------------\n")

	b.WriteString(
		lipgloss.JoinHorizontal(
			lipgloss.Center,
			lipgloss.NewStyle().Margin(0, 1).Render(f.centralKeyboardHalfView.View()),
			lipgloss.NewStyle().Margin(0, 1).Render(f.peripheralKeyboardHalfView.View()),

		),
	)

	return b.String()
}

func (f FlashView) Init() tea.Cmd {
	return NextStepCmd()
}

func ErrorCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return err
	}
}

func NextStepCmd() tea.Cmd {
	return func() tea.Msg {
		return NextStepMsg{}
	}
}

type NextStepMsg struct{}

func CopyFileCmd(identifier, src, dest string, dryRun bool) tea.Cmd {
	return func() tea.Msg {
		if dryRun {
			return FileCopiedMsg{Identifier: identifier}
		}
		err := files.CopyFile(src, dest)
		if err != nil {
			// Ignore input/output errors because they are likely due to the bootloader being removed after flashing
			if !strings.Contains(err.Error(), "input/output error") {
				return errors.Join(errors.New("error copying file from "+src+" to "+dest), err)
			}
		}
		return FileCopiedMsg{Identifier: identifier}
	}
}

type FileCopiedMsg struct {
	Identifier string
}
