package fetcher

import (
	"net/http"
	"slices"
)

type Headers struct {
	headers http.Header
}

func (headers *Headers) Get(header string) ([]string, bool) {
	values, ok := headers.headers[header]
	return values, ok
}

func (headers *Headers) GetFirst(header string) (string, bool) {
	values, ok := headers.headers[header]

	if len(values) == 0 {
		return "", false
	}

	return values[0], ok
}

func (headers *Headers) Each(consumer func(string, string)) {
	for header, values := range headers.headers {
		for _, value := range values {
			consumer(header, value)
		}
	}
}

func (headers *Headers) EachExcluding(excludes []string, consumer func(string, string)) {
	for header, values := range headers.headers {
		if slices.Contains(excludes, header) {
			continue
		}

		for _, value := range values {
			consumer(header, value)
		}
	}
}
