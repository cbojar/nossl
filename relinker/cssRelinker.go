package relinker

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var cssURLsWithDoubleQuotes *regexp.Regexp = regexp.MustCompile(`url\(\s*"[^"]+"\s*\)`)
var cssURLsWithSingleQuotes *regexp.Regexp = regexp.MustCompile(`url\(\s*'[^']+'\s*\)`)

type CSSRelinker struct {
	proxier *URLProxier
}

func NewCSSRelinker(proxier *URLProxier) *CSSRelinker {
	return &CSSRelinker{proxier: proxier}
}

func (relinker *CSSRelinker) Relink(proxiedURL *url.URL, css string) string {
	relinked := cssURLsWithDoubleQuotes.ReplaceAllStringFunc(css, relinker.relinkURL(proxiedURL, '"'))
	relinked = cssURLsWithSingleQuotes.ReplaceAllStringFunc(relinked, relinker.relinkURL(proxiedURL, '\''))

	return relinked
}

func (relinker *CSSRelinker) relinkURL(proxiedURL *url.URL, quote rune) func(string) string {
	return func(match string) string {
		if urlString := parseCSSURL(match); strings.HasPrefix(urlString, "data:") {
			return match
		} else if valueURL, err := url.Parse(urlString); err != nil {
			return match
		} else {
			return fmt.Sprintf(
				"url(%[1]c%[2]s%[1]c)",
				quote, relinker.proxier.Proxy(proxiedURL, valueURL))
		}
	}
}

func parseCSSURL(match string) string {
	urlString := strings.TrimPrefix(match, "url(")
	urlString = strings.TrimSuffix(urlString, ")")
	urlString = strings.TrimSpace(urlString)
	urlString = urlString[1 : len(urlString)-1]

	return urlString
}
