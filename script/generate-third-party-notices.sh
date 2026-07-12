#!/usr/bin/env bash

set -euo pipefail

output="${1:-THIRD_PARTY_NOTICES.txt}"
repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
work_dir="$(mktemp -d)"
trap 'rm -rf "$work_dir"' EXIT

targets=(
  darwin-amd64
  darwin-arm64
  freebsd-386
  freebsd-amd64
  freebsd-arm64
  linux-386
  linux-amd64
  linux-arm
  linux-arm64
  windows-386
  windows-amd64
  windows-arm64
)

cd "$repo_root"

for target in "${targets[@]}"; do
  goos="${target%-*}"
  goarch="${target#*-}"
  CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" \
    go list -deps -f '{{with .Module}}{{.Path}}|{{.Version}}|{{.Dir}}{{end}}' \
    ./cmd/gh-attach >> "$work_dir/modules"
done

sort -u "$work_dir/modules" | sed '/^[[:space:]]*$/d' > "$work_dir/modules.sorted"

{
  echo "Third-party license notices for gh-attach"
  echo
  echo "This file covers modules linked into at least one published platform binary."

  while IFS='|' read -r module version directory; do
    if [[ "$module" == "github.com/harley/gh-attach" ]]; then
      continue
    fi

    license_file_list="$work_dir/license-files"
    find "$directory" -maxdepth 2 -type f \
      \( -iname 'LICENSE*' -o -iname 'COPYING*' -o -iname 'NOTICE*' \) \
      | sort > "$license_file_list"

    if [[ ! -s "$license_file_list" ]]; then
      echo "error: no license file found for $module $version" >&2
      exit 1
    fi

    echo
    echo "================================================================================"
    echo "$module $version"

    while IFS= read -r license_file; do
      relative_path="${license_file#"$directory"/}"
      echo
      echo "--- $relative_path ---"
      cat "$license_file"
    done < "$license_file_list"
  done < "$work_dir/modules.sorted"
} > "$output"
