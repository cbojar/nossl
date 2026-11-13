package handlers

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/cbojar/nossl/fetcher"
	"github.com/cbojar/nossl/relinker"
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
	if toProxy, err := parseProxyURLQuery(context); err != nil {
		context.String(400, "%s", err)
	} else if response, err := handler.fetcher.Fetch("GET", toProxy, context.Request.Header); err != nil {
		context.String(400, "%s", err)
	} else {
		proxier := relinker.NewURLProxier(buildProxyBase(context))
		responseRelinker := relinker.RelinkerFor(response.MimeType(), proxier)
		responseBody := responseRelinker.Relink(toProxy, response.BodyString())

		applyProxiedHeaders(response, context)
		context.String(response.StatusCode(), "%s", responseBody)
	}
}

func parseProxyURLQuery(context *gin.Context) (*url.URL, error) {
	if queryURL, exists := context.GetQuery("url"); !exists {
		return nil, fmt.Errorf("No URL supplied")
	} else {
		return parseProxyURL(queryURL)
	}
}

func parseProxyURL(urlString string) (*url.URL, error) {
	if !hasProtocol(urlString) {
		return parseProxyURL(fmt.Sprintf("https://%s", urlString))
	} else {
		return url.Parse(urlString)
	}
}

func hasProtocol(urlString string) bool {
	return strings.HasPrefix(urlString, "http://") || strings.HasPrefix(urlString, "https://")
}

func buildProxyBase(context *gin.Context) *url.URL {
	proxyBase, _ := url.Parse(fmt.Sprintf("http://%s/proxy", context.Request.Host))
	return proxyBase
}

func applyProxiedHeaders(response *fetcher.Response, context *gin.Context) {
	response.Headers().EachExcluding(
		[]string{"Content-Length"},
		context.Writer.Header().Add)

	context.Writer.Header().Add("Via", "nossl 1.0")
}
