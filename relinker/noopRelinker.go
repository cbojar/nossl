package relinker

import "net/url"

type noopRelinker struct{}

func NewNoopRelinker() *noopRelinker {
	return &noopRelinker{}
}

func (*noopRelinker) Relink(proxiedURL *url.URL, content string) string {
	return content
}
