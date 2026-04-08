package app

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/harley/gh-attach/internal/render"
)

type prCreateOptions struct {
	Repo        string
	AttachPaths []string
	Help        bool
	Forwarded   []string
}

func runPRCreate(args []string) error {
	opts, err := parsePRCreateOptions(args)
	if err != nil {
		return err
	}
	if opts.Help {
		fmt.Print(usage)
		return nil
	}
	if len(opts.AttachPaths) == 0 {
		return fmt.Errorf("pr create requires at least one --attach file")
	}

	assets, err := uploadFiles(opts.Repo, opts.AttachPaths)
	if err != nil {
		return err
	}

	lines := make([]string, 0, len(assets))
	for _, asset := range assets {
		lines = append(lines, render.Markdown(asset))
	}

	baseBody, forwarded, err := extractPRBody(opts.Forwarded)
	if err != nil {
		return err
	}

	body := joinBody(baseBody, lines)
	cmdArgs := append([]string{"pr", "create"}, forwarded...)
	cmdArgs = append(cmdArgs, "--body", body)

	cmd := exec.Command("gh", cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("gh pr create: %w", err)
	}

	return nil
}

func parsePRCreateOptions(args []string) (*prCreateOptions, error) {
	opts := &prCreateOptions{}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "-h" || arg == "--help":
			opts.Help = true
		case arg == "--attach":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("--attach requires a file path")
			}
			i++
			opts.AttachPaths = append(opts.AttachPaths, args[i])
		case len(arg) > len("--attach=") && arg[:9] == "--attach=":
			opts.AttachPaths = append(opts.AttachPaths, arg[len("--attach="):])
		case arg == "--repo":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("--repo requires owner/repo")
			}
			i++
			opts.Repo = args[i]
		case len(arg) > len("--repo=") && arg[:7] == "--repo=":
			opts.Repo = arg[len("--repo="):]
		default:
			opts.Forwarded = append(opts.Forwarded, arg)
		}
	}

	return opts, nil
}

func extractPRBody(args []string) (string, []string, error) {
	var body string
	out := make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--body":
			if i+1 >= len(args) {
				return "", nil, fmt.Errorf("--body requires a value")
			}
			i++
			body = args[i]
		case len(arg) > len("--body=") && arg[:7] == "--body=":
			body = arg[len("--body="):]
		case arg == "--body-file":
			if i+1 >= len(args) {
				return "", nil, fmt.Errorf("--body-file requires a path")
			}
			i++
			content, err := os.ReadFile(args[i])
			if err != nil {
				return "", nil, fmt.Errorf("read --body-file %s: %w", args[i], err)
			}
			body = string(content)
		case len(arg) > len("--body-file=") && arg[:12] == "--body-file=":
			path := arg[len("--body-file="):]
			content, err := os.ReadFile(path)
			if err != nil {
				return "", nil, fmt.Errorf("read --body-file %s: %w", path, err)
			}
			body = string(content)
		case arg == "--fill" || arg == "--fill-first" || arg == "--fill-verbose":
			return "", nil, fmt.Errorf("%s is not supported with --attach yet; pass --body or --body-file instead", arg)
		default:
			out = append(out, arg)
		}
	}

	return body, out, nil
}
