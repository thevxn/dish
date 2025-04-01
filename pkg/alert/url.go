package alert

import (
	"fmt"
	"net/url"
	"slices"
	"strings"
)

var defaultSchemes = []string{"http", "https"}

// parseAndValidateURL parses and validates a URL with strict scheme requirements.
// The supportedSchemes parameter allows customizing allowed protocols (defaults to http/https if nil).
func parseAndValidateURL(rawURL string, supportedSchemes []string) (*url.URL, error) {
	if strings.TrimSpace(rawURL) == "" {
		return nil, fmt.Errorf("URL cannot be empty")
	}

	if supportedSchemes == nil {
		supportedSchemes = defaultSchemes
	}

	// Parse the provided URL
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %w", err)
	}

	// Validate the parsed URL
	switch {
	case parsedURL.Scheme == "":
		return nil, fmt.Errorf("protocol must be specified in the provided URL (e.g. https://...)")

	case !slices.Contains(supportedSchemes, parsedURL.Scheme):
		return nil, fmt.Errorf("unsupported protocol provided in URL: %s (supported protocols: %v)", parsedURL.Scheme, supportedSchemes)

	case parsedURL.Host == "":
		return nil, fmt.Errorf("URL must contain a host")
	}

	return parsedURL, nil
}
