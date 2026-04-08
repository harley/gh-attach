package repo

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type Info struct {
	Owner string
	Name  string
	ID    int
}

var (
	sshRemoteRe   = regexp.MustCompile(`git@github\.com:([^/]+)/([^/]+?)(?:\.git)?$`)
	httpsRemoteRe = regexp.MustCompile(`https://github\.com/([^/]+)/([^/]+?)(?:\.git)?$`)
)

func Resolve(repoArg string) (*Info, error) {
	var owner string
	var name string

	if repoArg != "" {
		var err error
		owner, name, err = split(repoArg)
		if err != nil {
			return nil, err
		}
	} else {
		var err error
		owner, name, err = fromRemote()
		if err != nil {
			return nil, err
		}
	}

	id, err := lookupID(owner, name)
	if err != nil {
		return nil, err
	}

	return &Info{Owner: owner, Name: name, ID: id}, nil
}

func split(input string) (string, string, error) {
	parts := strings.SplitN(input, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid repo %q, expected owner/repo", input)
	}
	return parts[0], parts[1], nil
}

func fromRemote() (string, string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	out, err := cmd.Output()
	if err != nil {
		return "", "", fmt.Errorf("could not infer repo from git remote; pass --repo owner/repo")
	}

	remote := strings.TrimSpace(string(out))
	if match := sshRemoteRe.FindStringSubmatch(remote); match != nil {
		return match[1], match[2], nil
	}
	if match := httpsRemoteRe.FindStringSubmatch(remote); match != nil {
		return match[1], match[2], nil
	}

	return "", "", fmt.Errorf("unsupported remote URL %q", remote)
}

func lookupID(owner, name string) (int, error) {
	cmd := exec.Command("gh", "api", fmt.Sprintf("repos/%s/%s", owner, name), "--jq", ".id")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	out, err := cmd.Output()
	if err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return 0, fmt.Errorf("lookup repo id for %s/%s: %s", owner, name, msg)
	}

	id, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		return 0, fmt.Errorf("parse repo id for %s/%s: %w", owner, name, err)
	}
	return id, nil
}
