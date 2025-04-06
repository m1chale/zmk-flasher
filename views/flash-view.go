package views

import (
	"errors"
	"strconv"
	"strings"
	"zmk-flasher/files"

	tea "github.com/charmbracelet/bubbletea"
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
	step                flashStep
	leftBootloaderFile  string
	rightBootloaderFile string

	mountLeftBootloaderView  MountBootloaderView
	mountRightBootloaderView MountBootloaderView

	dryRun bool

	err error
}

func NewFlashView(leftBootloaderFile, rightBootloaderFile, leftMountPoint, rightMountPoint string, dryRun bool) FlashView {
	return FlashView{
		step:                MountLeftBootloader - 1,
		leftBootloaderFile:  leftBootloaderFile,
		rightBootloaderFile: rightBootloaderFile,
		mountLeftBootloaderView: MountBootloaderView{
			bootloaderName: "left",
			mountPath:      leftMountPoint,
		},
		mountRightBootloaderView: MountBootloaderView{
			bootloaderName: "right",
			mountPath:      rightMountPoint,
		},
		dryRun: dryRun,
	}
}

func (f FlashView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	commands := []tea.Cmd{}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" {
			return f, tea.Quit
		}
	case NextStepMsg:
		f.step++
		if f.step == Done {
			return f, tea.Quit
		}
	case FileCopiedMsg:
		if msg.Identifier == "left" {
			return f, NextStepCmd()
		}
		if msg.Identifier == "right" {
			return f, NextStepCmd()
		}
	case error:
		println("error")
		println(msg.Error())
		return f, tea.Quit
	}

	switch f.step {
	case MountLeftBootloader:
		model, cmd := f.mountLeftBootloaderView.Update(msg)
		f.mountLeftBootloaderView = model.(MountBootloaderView)
		if cmd != nil {
			commands = append(commands, cmd)
		}
	case MountRightBootloader:
		model, cmd := f.mountRightBootloaderView.Update(msg)
		f.mountRightBootloaderView = model.(MountBootloaderView)
		if cmd != nil {
			commands = append(commands, cmd)
		}
	case FlashLeftBootloader:
		commands = append(commands, CopyFileCmd("left", f.leftBootloaderFile, f.mountLeftBootloaderView.mountPath+"/file.uf2", f.dryRun))
	case FlashRightBootloader:
		commands = append(commands, CopyFileCmd("right", f.rightBootloaderFile, f.mountRightBootloaderView.mountPath+"/file.uf2", f.dryRun))
	}

	return f, tea.Batch(commands...)
}

func (f FlashView) View() string {
	b := strings.Builder{}
	b.WriteString("Flash View\n")
	b.WriteString("Left Bootloader File: " + f.leftBootloaderFile + "\n")
	b.WriteString("Right Bootloader File: " + f.rightBootloaderFile + "\n")
	b.WriteString("Left Mount Point: " + f.mountLeftBootloaderView.mountPath + "\n")
	b.WriteString("Right Mount Point: " + f.mountRightBootloaderView.mountPath + "\n")
	b.WriteString("Dry Run: " + strconv.FormatBool(f.dryRun) + "\n")
	if f.err != nil {
		b.WriteString("Error: " + f.err.Error() + "\n")
	}
	b.WriteString("----------------\n")

	switch f.step {
	case MountLeftBootloader:
		b.WriteString(f.mountLeftBootloaderView.View())
	case MountRightBootloader:
		b.WriteString(f.mountRightBootloaderView.View())
	case FlashLeftBootloader:
		b.WriteString("Flashing left bootloader\n")
	case FlashRightBootloader:
		b.WriteString("Flashing right bootloader\n")
	case Done:
		b.WriteString("Done\n")
	}

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
