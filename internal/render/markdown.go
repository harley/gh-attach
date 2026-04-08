package render

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/harley/gh-attach/internal/attach"
)

func Markdown(asset attach.Asset) string {
	if isImage(asset) {
		name := strings.TrimSuffix(asset.Filename, filepath.Ext(asset.Filename))
		return fmt.Sprintf("![%s](%s)", name, asset.URL)
	}
	return fmt.Sprintf("[%s](%s)", asset.Filename, asset.URL)
}

func isImage(asset attach.Asset) bool {
	if strings.HasPrefix(asset.ContentType, "image/") {
		return true
	}

	switch strings.ToLower(filepath.Ext(asset.Filename)) {
	case ".png", ".jpg", ".jpeg", ".gif", ".svg", ".webp", ".bmp", ".ico", ".tif", ".tiff", ".avif":
		return true
	default:
		return false
	}
}
