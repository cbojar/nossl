package fetcher

import (
	"io"
	"net/http"
	"net/url"
)

type Fetcher struct {
	client *http.Client
}

func NewFetcher() *Fetcher {
	return &Fetcher{
		client: http.DefaultClient,
	}
}

func (fetcher *Fetcher) Fetch(method string, url *url.URL, header http.Header) (int, http.Header, string, error) {
	header.Del("Accept-Encoding") // Prevent getting a compressed response back

	request := &http.Request{
		Method: method,
		URL:    url,
		Header: header,
	}

	response, err := fetcher.client.Do(request)
	if err != nil {
		return -1, nil, "", err
	}

	defer response.Body.Close()

	if body, err := io.ReadAll(response.Body); err != nil {
		return response.StatusCode, response.Header, "", err
	} else {
		return response.StatusCode, response.Header, string(body), nil
	}
}
