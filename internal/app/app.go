package app

import (
	"fmt"
	"strings"
)

const usage = `gh-attach

Usage:
  gh attach upload [--repo owner/repo] [--json] <file>...
  gh attach pr create [--repo owner/repo] --attach <file>... [gh pr create flags...]

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
	case len(args) >= 2 && args[0] == "pr" && args[1] == "create":
		return runPRCreate(args[2:])
	default:
		return runUpload(args)
	}
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
