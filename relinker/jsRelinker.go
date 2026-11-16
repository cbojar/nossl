package relinker

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var valuesWithDoubleQuotes *regexp.Regexp = regexp.MustCompile(`"https?://[^"]+"`)
var valuesWithSingleQuotes *regexp.Regexp = regexp.MustCompile(`'https?://[^']+'`)
var valuesWithTickQuotes *regexp.Regexp = regexp.MustCompile("`https?://[^`]+`")

type JSRelinker struct {
	proxier *URLProxier
}

func NewJSRelinker(proxier *URLProxier) *JSRelinker {
	return &JSRelinker{proxier: proxier}
}

func (relinker *JSRelinker) Relink(proxiedURL *url.URL, html string) string {
	relinked := valuesWithDoubleQuotes.ReplaceAllStringFunc(html, relinker.relinkValue(proxiedURL, '"'))
	relinked = valuesWithSingleQuotes.ReplaceAllStringFunc(relinked, relinker.relinkValue(proxiedURL, '\''))
	relinked = valuesWithTickQuotes.ReplaceAllStringFunc(relinked, relinker.relinkValue(proxiedURL, '`'))

	return relinked
}

func (relinker *JSRelinker) relinkValue(proxiedURL *url.URL, quote rune) func(string) string {
	return func(match string) string {
		if value := parseValue(match); strings.HasPrefix(value, "data:") {
			return match
		} else if valueURL, err := url.Parse(value); err != nil {
			return match
		} else {
			return fmt.Sprintf(
				"%[1]c%[2]s%[1]c",
				quote, relinker.proxier.Proxy(proxiedURL, valueURL))
		}
	}
}

func parseValue(match string) string {
	return match[1 : len(match)-1]
}
