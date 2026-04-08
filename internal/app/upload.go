package app

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/harley/gh-attach/internal/attach"
	"github.com/harley/gh-attach/internal/browser"
	"github.com/harley/gh-attach/internal/render"
	"github.com/harley/gh-attach/internal/repo"
)

type uploadOptions struct {
	Repo  string
	JSON  bool
	Help  bool
	Files []string
}

func runUpload(args []string) error {
	opts, err := parseUploadOptions(args)
	if err != nil {
		return err
	}
	if opts.Help {
		fmt.Print(usage)
		return nil
	}
	if len(opts.Files) == 0 {
		return fmt.Errorf("no files provided")
	}

	assets, err := uploadFiles(opts.Repo, opts.Files)
	if err != nil {
		return err
	}

	if opts.JSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(assets)
	}

	for _, asset := range assets {
		fmt.Println(render.Markdown(asset))
	}
	return nil
}

func parseUploadOptions(args []string) (*uploadOptions, error) {
	opts := &uploadOptions{}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "-h" || arg == "--help":
			opts.Help = true
		case arg == "--json":
			opts.JSON = true
		case arg == "--repo":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("--repo requires owner/repo")
			}
			i++
			opts.Repo = args[i]
		case len(arg) > len("--repo=") && arg[:7] == "--repo=":
			opts.Repo = arg[len("--repo="):]
		case len(arg) > 0 && arg[0] == '-':
			return nil, fmt.Errorf("unknown flag: %s", arg)
		default:
			opts.Files = append(opts.Files, arg)
		}
	}

	return opts, nil
}

func uploadFiles(repoArg string, files []string) ([]attach.Asset, error) {
	repoInfo, err := repo.Resolve(repoArg)
	if err != nil {
		return nil, err
	}

	cookie, err := browser.GetGitHubSession()
	if err != nil {
		return nil, err
	}

	client := attach.NewClient(cookie)
	results := make([]attach.Asset, 0, len(files))
	for _, path := range files {
		asset, err := attach.Upload(client, repoInfo.Owner, repoInfo.Name, repoInfo.ID, path)
		if err != nil {
			return nil, fmt.Errorf("upload %s: %w", path, err)
		}
		results = append(results, *asset)
	}

	return results, nil
}
