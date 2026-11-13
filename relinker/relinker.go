package relinker

import "net/url"

type Relinker interface {
	Relink(proxiedURL *url.URL, content string) string
}
