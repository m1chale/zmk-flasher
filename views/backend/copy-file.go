package backend

import (
	"errors"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/new-er/zmk-flasher/files"
)

func CopyFileCmd(src, dest string, dryRun bool) tea.Cmd {
	return func() tea.Msg {
		if dryRun {
			return FileCopiedMsg{}
		}
		err := files.CopyFile(src, dest)
		if err != nil {
			// Ignore input/output errors because they are likely due to the bootloader being removed after flashing
			if !strings.Contains(err.Error(), "input/output error") {
				return errors.Join(errors.New("error copying file from "+src+" to "+dest), err)
			}
		}
		return FileCopiedMsg{}
	}
}

type FileCopiedMsg struct {
}
