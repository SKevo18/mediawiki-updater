package httputil

import (
	"fmt"
	"net/http"
)

const UserAgent = "mediawiki-updater (https://github.com/SKevo18/mediawiki-updater)"

// Get performs an HTTP GET request with a proper User-Agent header.
func Get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", UserAgent)
	return http.DefaultClient.Do(req)
}
