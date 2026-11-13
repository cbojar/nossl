package relinker

import (
	"fmt"
	"html"
	"net/url"
	"regexp"
	"strings"
)

var attributesWithDoubleQuotes *regexp.Regexp = regexp.MustCompile(`(src|href|action|content)\s*=\s*("[^"]+")`)
var attributesWithSingleQuotes *regexp.Regexp = regexp.MustCompile(`(src|href|action|content)\s*=\s*('[^']+')`)

type HTMLRelinker struct {
	proxier *URLProxier
}

func NewHTMLRelinker(proxier *URLProxier) *HTMLRelinker {
	return &HTMLRelinker{proxier: proxier}
}

func (relinker *HTMLRelinker) Relink(proxiedURL *url.URL, html string) string {
	relinked := attributesWithDoubleQuotes.ReplaceAllStringFunc(html, relinker.relinkAttribute(proxiedURL, '"'))
	relinked = attributesWithSingleQuotes.ReplaceAllStringFunc(relinked, relinker.relinkAttribute(proxiedURL, '\''))

	return relinked
}

func (relinker *HTMLRelinker) relinkAttribute(proxiedURL *url.URL, quote rune) func(string) string {
	return func(match string) string {
		if attribute, value := parseAttribute(match); strings.HasPrefix(value, "data:") {
			return match
		} else if valueURL, err := url.Parse(value); err != nil {
			return match
		} else {
			return fmt.Sprintf(
				"%[2]s=%[1]c%[3]s%[1]c",
				quote, attribute, relinker.proxy(proxiedURL, valueURL))
		}
	}
}

func (relinker *HTMLRelinker) proxy(proxiedURL *url.URL, valueURL *url.URL) string {
	return html.EscapeString(relinker.proxier.Proxy(proxiedURL, valueURL).String())
}

func parseAttribute(match string) (string, string) {
	attribute, value, _ := strings.Cut(match, "=")

	attribute = strings.TrimSpace(attribute)

	value = strings.TrimSpace(value)
	value = value[1 : len(value)-1]
	value = html.UnescapeString(value)

	return attribute, value
}
