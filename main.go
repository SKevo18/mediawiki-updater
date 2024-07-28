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

	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		fmt.Println("Error creating target directory:", err)
		return
	}

	if err := downloadMediaWikiCore(targetDir); err != nil {
		fmt.Println("Error downloading MediaWiki core:", err)
		return
	}

	if err := downloadExtensionsAndSkins(targetDir); err != nil {
		fmt.Println("Error downloading extensions and skins:", err)
		return
	}

	fmt.Println("MediaWiki core and extensions/skins downloaded successfully.")
}

// Downloads the MediaWiki core from the specified URL and extracts it into the specified path.
func downloadMediaWikiCore(targetDir string) error {
	return downloadAndExtract(mediaWikiURL, targetDir, func(path string) string {
		parts := strings.Split(path, string(filepath.Separator))
		parts = parts[1:]
		return strings.Join(parts, string(filepath.Separator))
	})
}

// Downloads extensions and skins for MediaWiki from specified `targetDir`.
func downloadExtensionsAndSkins(targetDir string) error {
	if err := downloadFromFile("extensions.txt", filepath.Join(targetDir, "extensions")); err != nil {
		return fmt.Errorf("error downloading extensions: %w", err)
	}
	if err := downloadFromFile("skins.txt", filepath.Join(targetDir, "skins")); err != nil {
		return fmt.Errorf("error downloading skins: %w", err)
	}

	return nil
}

// Downloads extensions or skins from the specified text file path (one extension/skin name per line).
// If the `filePath` does not exist, the function returns without error (nothing happens).
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

// Downloads the specified extension or skin from the specified URL and extracts it into the specified path.
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

// Downloads the specified ".tar.gz" archive file from the specified URL and extracts it into the specified path.
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

// Extracts the contents of the specified ".tar.gz" archive file into the specified path.
func extractGzArchive(file io.Reader, targetPath string, renamer extract.Renamer) error {
	if err := extract.Gz(context.TODO(), file, targetPath, renamer); err != nil {
		return err
	}

	return nil
}

// Downloads the file from the specified URL and saves it as `targetPath`.
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

// Obtains names of all extensions from the specified text file path
// (one extension name per line) and returns them as a slice of strings.
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

// Obtains a download URL for an extension or skin for specified MediaWiki extension from the specified
// extension dist URL (e. g.: https://extdist.wmflabs.org/dist/extensions/).
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
