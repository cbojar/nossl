package handlers

import (
	"fmt"
	"maps"
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
	requestQuery := context.Request.URL.Query().Clone()
	if queryURL := context.Query("url"); len(queryURL) == 0 {
		return nil, fmt.Errorf("No URL supplied")
	} else if strings.HasPrefix(queryURL, "http://") {
		return buildURL("http", strings.TrimPrefix(queryURL, "http://"), requestQuery)
	} else if strings.HasPrefix(queryURL, "https://") {
		return buildURL("https", strings.TrimPrefix(queryURL, "https://"), requestQuery)
	} else if queryProto := context.Query("proto"); len(queryProto) > 0 {
		return buildURL(queryProto, queryURL, requestQuery)
	} else {
		return buildURL("https", queryURL, requestQuery)
	}
}

func buildURL(proto string, urlString string, query url.Values) (*url.URL, error) {
	if parsed, err := url.Parse(fmt.Sprintf("%s://%s", proto, urlString)); err != nil {
		return nil, err
	} else {
		query.Del("url")
		query.Del("proto")
		parsedQuery := parsed.Query()
		maps.Insert(parsedQuery, maps.All(query))
		parsed.RawQuery = parsedQuery.Encode()
		return parsed, nil
	}
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
