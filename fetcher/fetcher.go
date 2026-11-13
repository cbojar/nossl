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

func (fetcher *Fetcher) Fetch(method string, url *url.URL, header http.Header) (*Response, error) {
	header.Del("Accept-Encoding") // Prevent getting a compressed response back

	request := &http.Request{
		Method: method,
		URL:    url,
		Header: header,
	}

	response, err := fetcher.client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	return &Response{
		statusCode: response.StatusCode,
		status:     response.Status,
		headers:    &Headers{headers: response.Header},
		body:       body,
	}, err
}
