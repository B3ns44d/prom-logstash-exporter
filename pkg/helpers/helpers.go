package helpers

import (
	"net/url"
	"strings"
)

func ParseURI(uri string) (*url.URL, error) {
	// Remove any trailing slashes from the URI
	if strings.HasSuffix(uri, "/") {
		uri = uri[0 : len(uri)-1]
	}

	// Parse the URI as a URL
	parsedURL, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	return parsedURL, nil
}
