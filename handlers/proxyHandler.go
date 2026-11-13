package handlers

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/cbojar/nossl/fetcher"
	"github.com/gin-gonic/gin"
)

type ProxyHandler struct {
	fetcher *fetcher.Fetcher
}

func NewProxyHandler(fetcher *fetcher.Fetcher) *ProxyHandler {
	return &ProxyHandler{
		fetcher: fetcher,
	}
}

func (handler *ProxyHandler) Handle(context *gin.Context) {
	if toProxy, err := parseProxyURL(context); err != nil {
		context.String(400, "%s", err)
	} else if response, err := handler.fetcher.Fetch("GET", toProxy, context.Request.Header); err != nil {
		context.String(400, "%s", err)
	} else {
		applyProxiedHeaders(response, context)

		proxyBase, _ := url.Parse(fmt.Sprintf("http://%s/proxy", context.Request.Host))
		responseBody := proxyLinkURLs(response.BodyString(), toProxy, proxyBase)

		context.String(response.StatusCode(), "%s", responseBody)
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

func applyProxiedHeaders(response *fetcher.Response, context *gin.Context) {
	response.Headers().EachExcluding(
		[]string{"Content-Length"},
		context.Writer.Header().Add)

	context.Writer.Header().Add("Via", "nossl 1.0")
}
