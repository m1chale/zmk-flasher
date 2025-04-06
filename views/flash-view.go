package views

import (
	"errors"
	"strconv"
	"strings"
	"github.com/new-er/zmk-flasher/files"

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
	rightBootloaderFile string
	leftBootloaderFile  string

	mountRightBootloaderView MountBootloaderView
	mountLeftBootloaderView  MountBootloaderView

	dryRun bool
}

func NewFlashView(leftBootloaderFile, rightBootloaderFile, leftMountPoint, rightMountPoint string, dryRun bool) FlashView {
	return FlashView{
		step:                Init,
		rightBootloaderFile: rightBootloaderFile,
		leftBootloaderFile:  leftBootloaderFile,
		mountRightBootloaderView: MountBootloaderView{
			bootloaderName: "right",
			mountPath:      rightMountPoint,
		},
		mountLeftBootloaderView: MountBootloaderView{
			bootloaderName: "left",
			mountPath:      leftMountPoint,
		},
		dryRun: dryRun,
	}
}

func (f FlashView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if f.step > MountRightBootloader && f.step < FlashRightBootloader {
		if err := f.mountRightBootloaderView.EnsureMountPathExists(); err != nil {
			println("right bootloader was removed unexpectedly")
			return f, tea.Quit
		}
	}
	if f.step > MountLeftBootloader && f.step < FlashLeftBootloader {
		if err := f.mountLeftBootloaderView.EnsureMountPathExists(); err != nil {
			println("left bootloader was removed unexpectedly")
			return f, tea.Quit
		}
	}
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
	case MountRightBootloader:
		model, cmd := f.mountRightBootloaderView.Update(msg)
		f.mountRightBootloaderView = model.(MountBootloaderView)
		if cmd != nil {
			commands = append(commands, cmd)
		}
	case MountLeftBootloader:
		model, cmd := f.mountLeftBootloaderView.Update(msg)
		f.mountLeftBootloaderView = model.(MountBootloaderView)
		if cmd != nil {
			commands = append(commands, cmd)
		}
	case FlashRightBootloader:
		commands = append(commands, CopyFileCmd("right", f.rightBootloaderFile, f.mountRightBootloaderView.mountPath+"/file.uf2", f.dryRun))
	case FlashLeftBootloader:
		commands = append(commands, CopyFileCmd("left", f.leftBootloaderFile, f.mountLeftBootloaderView.mountPath+"/file.uf2", f.dryRun))
	}

	return f, tea.Batch(commands...)
}

func (f FlashView) View() string {
	b := strings.Builder{}
	b.WriteString("Zmk Flasher\n")
	b.WriteString("Right Bootloader File: " + f.rightBootloaderFile + "\n")
	b.WriteString("Left Bootloader File: " + f.leftBootloaderFile + "\n")
	b.WriteString("Right Mount Point: " + f.mountRightBootloaderView.mountPath + "\n")
	b.WriteString("Left Mount Point: " + f.mountLeftBootloaderView.mountPath + "\n")
	b.WriteString("Dry Run: " + strconv.FormatBool(f.dryRun) + "\n")
	b.WriteString("Press 'q' to quit\n")
	b.WriteString("----------------\n")

	switch f.step {
	case MountRightBootloader:
		b.WriteString(f.mountRightBootloaderView.View())
	case MountLeftBootloader:
		b.WriteString(f.mountLeftBootloaderView.View())
	case FlashRightBootloader:
		b.WriteString("Flashing right bootloader\n")
	case FlashLeftBootloader:
		b.WriteString("Flashing left bootloader\n")
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
