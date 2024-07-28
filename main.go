package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/codeclysm/extract/v3"
)

const (
	versionTag   string = "REL1_42"
	mediaWikiURL string = "https://releases.wikimedia.org/mediawiki/1.42/mediawiki-core-1.42.1.tar.gz"
	extensionURL string = "https://extdist.wmflabs.org/dist/extensions/"
	skinURL      string = "https://extdist.wmflabs.org/dist/skins/"
)

func main() {
	targetDir := "."
	if len(os.Args) > 1 {
		targetDir = os.Args[1]
	}

	tempDir, err := os.MkdirTemp("", "mediawiki-temp-")
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		return
	}
	defer os.RemoveAll(tempDir)

	if err := downloadMediaWikiCore(tempDir); err != nil {
		fmt.Println("Error downloading MediaWiki core:", err)
		return
	}

	if err := downloadExtensionsAndSkins(tempDir); err != nil {
		fmt.Println("Error downloading extensions and skins:", err)
		return
	}

	ignorePaths := []string{"LocalSettings.php", ".htaccess", "images"}
	if err := copyContents(tempDir, targetDir, ignorePaths); err != nil {
		fmt.Println("Error copying contents to target directory:", err)
		return
	}

	fmt.Println("MediaWiki core and extensions/skins downloaded and copied successfully.")
}

func copyContents(src, dst string, ignorePaths []string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		for _, ignorePath := range ignorePaths {
			if relPath == ignorePath || strings.HasPrefix(relPath, ignorePath+string(os.PathSeparator)) {
				return nil
			}
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		if err != nil {
			return err
		}

		return os.Chmod(dstPath, info.Mode())
	})
}

func downloadMediaWikiCore(targetDir string) error {
	return downloadAndExtract(mediaWikiURL, targetDir, func(path string) string {
		parts := strings.Split(path, string(filepath.Separator))
		parts = parts[1:]
		return strings.Join(parts, string(filepath.Separator))
	})
}

func downloadExtensionsAndSkins(targetDir string) error {
	if err := downloadFromFile("extensions.txt", filepath.Join(targetDir, "extensions")); err != nil {
		return fmt.Errorf("error downloading extensions: %w", err)
	}
	if err := downloadFromFile("skins.txt", filepath.Join(targetDir, "skins")); err != nil {
		return fmt.Errorf("error downloading skins: %w", err)
	}

	return nil
}

func downloadFromFile(filePath string, targetDir string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	}

	extensions, err := getExtensions(filePath)
	if err != nil {
		return err
	}

	for _, extension := range extensions {
		if err := downloadDist(extension, versionTag, targetDir); err != nil {
			return err
		}
	}

	return nil
}

func downloadDist(extensionName string, version string, targetDir string) error {
	downloadUrl, err := getDownloadUrl(extensionURL, extensionName, version)
	if err != nil {
		return err
	}
	if downloadUrl == "" {
		return fmt.Errorf("extension not found: %s", extensionName)
	}

	return downloadAndExtract(downloadUrl, targetDir, nil)
}

func downloadAndExtract(url string, targetDir string, renamer extract.Renamer) error {
	tempFile, err := os.CreateTemp("", "mw-temp-*.tar.gz")
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())

	if err := downloadFile(url, tempFile.Name()); err != nil {
		return err
	}
	if err := extractGzArchive(tempFile, targetDir, renamer); err != nil {
		return err
	}

	return nil
}

func extractGzArchive(file io.Reader, targetPath string, renamer extract.Renamer) error {
	if err := extract.Gz(context.TODO(), file, targetPath, renamer); err != nil {
		return err
	}

	return nil
}

func downloadFile(url string, targetPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = io.Copy(file, resp.Body); err != nil {
		return err
	}

	return nil
}

func getExtensions(path string) (extensions []string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		extensions = append(extensions, scanner.Text())
	}

	return extensions, nil
}

func getDownloadUrl(fromUrl string, extensionName string, version string) (downloadUrl string, err error) {
	resp, err := http.Get(fromUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), fmt.Sprintf("%s-%s", extensionName, version)) {
			downloadUrl = s.AttrOr("href", "")
		}
	})

	if downloadUrl == "" {
		return "", nil
	}
	return fromUrl + downloadUrl, nil
}
