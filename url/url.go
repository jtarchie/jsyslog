package url

import (
	"fmt"
	"net/url"
	"strings"
)

type URL struct {
	*url.URL
}

func (u *URL) String() string {
	if u.Scheme == "file" {
		return fmt.Sprintf("file://%s", u.Path)
	}

	return u.URL.String()
}

func Parse(rawURL string) (*URL, error) {
	if strings.HasPrefix(rawURL, "file://") {
		uri := &URL{
			URL: &url.URL{
				Scheme: "file",
				Path:   strings.TrimPrefix(rawURL, "file://"),
			},
		}
		return uri, nil
	}

	uri, err := url.Parse(rawURL)
	return &URL{
		URL: uri,
	}, err
}
