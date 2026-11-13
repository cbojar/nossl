package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

func ProxyHandler(context *gin.Context) {
	if toProxy, err := parseProxyURL(context); err != nil {
		context.String(400, "%s", err)
	} else if status, responseHeader, responseBody, err := fetchURL("GET", toProxy, context.Request.Header); err != nil {
		context.String(400, "%s", err)
	} else {
		applyProxiedHeaders(responseHeader, context)

		proxyBase, _ := url.Parse(fmt.Sprintf("http://%s/proxy", context.Request.Host))
		responseBody = proxyLinkURLs(responseBody, toProxy, proxyBase)

		context.String(status, "%s", responseBody)
	}
}

func parseProxyURL(context *gin.Context) (*url.URL, error) {
	toProxyString, exists := context.GetQuery("url")
	if !exists {
		return nil, fmt.Errorf("No URL supplied")
	}

	if !strings.HasPrefix(toProxyString, "http://") && !strings.HasPrefix(toProxyString, "https://") {
		toProxyString = fmt.Sprintf("https://%s", toProxyString)
	}

	toProxy, err := url.Parse(toProxyString)
	if err != nil {
		return nil, err
	}

	return toProxy, nil
}

func fetchURL(method string, url *url.URL, header http.Header) (int, http.Header, string, error) {
	header.Del("Accept-Encoding")
	request := &http.Request{
		Method: method,
		URL:    url,
		Header: header,
	}

	response, err := http.DefaultClient.Do(request)
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

func proxyLinkURLs(responseBody string, toProxy *url.URL, proxyBase *url.URL) string {
	rex := regexp.MustCompile(`(src|href|action)\s*=\s*("[^"]+"|'[^']+')`)
	responseBody = rex.ReplaceAllStringFunc(responseBody, func(match string) string {
		attribute, value := parseAttribute(match)
		valueURL, err := url.Parse(value)

		if err != nil {
			return match
		}

		resolved := toProxy.ResolveReference(valueURL)

		proxied := proxyBase.Clone()

		query := proxied.Query()
		query.Add("url", resolved.String())
		proxied.RawQuery = query.Encode()

		return fmt.Sprintf("%s=\"%s\"", attribute, proxied)
	})
	return responseBody
}

func parseAttribute(match string) (string, string) {
	split := strings.Split(match, "=")

	attribute := strings.TrimSpace(split[0])

	value := strings.TrimSpace(split[1])
	value = value[1 : len(value)-1]

	return attribute, value
}

func applyProxiedHeaders(responseHeader http.Header, context *gin.Context) {
	for header, values := range responseHeader {
		if header == "Content-Length" {
			continue
		}

		for _, value := range values {
			context.Writer.Header().Add(header, value)
		}
	}

	context.Writer.Header().Add("Via", "nossl 1.0")
}
