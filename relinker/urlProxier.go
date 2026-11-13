package relinker

import "net/url"

type URLProxier struct {
	proxyBase func() *url.URL
}

func NewURLProxier(proxyBase *url.URL) *URLProxier {
	return &URLProxier{proxyBase: proxyBase.Clone}
}

func (proxier *URLProxier) Proxy(proxiedURL *url.URL, unproxied *url.URL) *url.URL {
	proxied := proxier.proxyBase()
	resolved := proxiedURL.ResolveReference(unproxied)

	query := proxied.Query()
	query.Add("url", resolved.String())
	proxied.RawQuery = query.Encode()

	return proxied
}
