package relinker

import (
	"net/url"
)

type Relinker interface {
	Relink(proxiedURL *url.URL, content string) string
}

func RelinkerFor(mimeType string, proxier *URLProxier) Relinker {
	switch mimeType {
	case "text/css":
		return NewCSSRelinker(proxier)
	default:
		return NewHTMLRelinker(proxier)
	}
}
