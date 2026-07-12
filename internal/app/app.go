package app

import (
	"fmt"
	"strings"

	"github.com/harley/gh-attach/internal/browser"
)

const usage = `gh-attach

Usage:
  gh attach upload [--repo owner/repo] [--json] <file>...
  gh attach pr create [--repo owner/repo] --attach <file>... [gh pr create flags...]
  gh attach doctor

Shortcuts:
  gh attach <file>...                              Upload files and print markdown
  gh attach pr create --attach shot.png --title "Fix" --body "Context"

Notes:
  - Native uploads require an authenticated GitHub browser session.
  - Attachments render natively on GitHub and stay private on private repos.
  - The initial pr wrapper only supports attachment-aware body composition.
`

func Run(args []string) error {
	if len(args) == 0 {
		fmt.Print(usage)
		return nil
	}

	switch {
	case isHelp(args):
		fmt.Print(usage)
		return nil
	case args[0] == "upload":
		return runUpload(args[1:])
	case isDoctorCommand(args):
		fmt.Println(browser.Diagnose().String())
		return nil
	case len(args) >= 2 && args[0] == "pr" && args[1] == "create":
		return runPRCreate(args[2:])
	default:
		return runUpload(args)
	}
}

func isDoctorCommand(args []string) bool {
	return len(args) == 1 && args[0] == "doctor"
}

func isHelp(args []string) bool {
	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			return true
		}
	}
	return false
}

func joinBody(base string, lines []string) string {
	attachments := strings.Join(lines, "\n")
	if strings.TrimSpace(base) == "" {
		return attachments
	}
	if attachments == "" {
		return base
	}
	return strings.TrimRight(base, "\n") + "\n\n" + attachments
}
