package extractor

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/codeclysm/extract/v3"
)

// Extractor handles file extraction operations
type Extractor struct{}

// NewExtractor creates a new Extractor instance
func NewExtractor() *Extractor {
	return &Extractor{}
}

// ExtractArchive extracts a tar.gz archive to the specified directory
func (e *Extractor) ExtractArchive(reader io.Reader, targetDir string, renamer extract.Renamer) error {
	return extract.Gz(context.TODO(), reader, targetDir, renamer)
}

// ExtractMediaWikiCore extracts MediaWiki core with proper path manipulation
func (e *Extractor) ExtractMediaWikiCore(reader io.Reader, targetDir string) error {
	return e.ExtractArchive(reader, targetDir, func(path string) string {
		// Remove the first directory component (e.g., mediawiki-1.43.1/)
		parts := strings.Split(path, string(filepath.Separator))
		if len(parts) > 1 {
			parts = parts[1:]
		}
		return strings.Join(parts, string(filepath.Separator))
	})
}

// CopyContents copies files from source to destination, ignoring specified paths
func (e *Extractor) CopyContents(src, dst string, ignorePaths []string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		// Check if path should be ignored
		for _, ignorePath := range ignorePaths {
			if relPath == ignorePath || strings.HasPrefix(relPath, ignorePath+string(os.PathSeparator)) {
				return nil
			}
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		return e.copyFile(path, dstPath, info.Mode())
	})
}

// copyFile copies a single file from source to destination
func (e *Extractor) copyFile(src, dst string, mode os.FileMode) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return os.Chmod(dst, mode)
}
