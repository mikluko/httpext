package httpext

import (
	"crypto/tls"
	"net/http"
	"net/url"

	"github.com/mikluko/httpext/tlshuffle"
	"github.com/mikluko/httpext/useragent"
)

type Option func(*Client)

func (opt Option) apply(c *Client) {
	opt(c)
}

func WithProxyURL(proxy *url.URL) Option {
	return func(c *Client) {
		if c.Transport == nil {
			return
		}
		t, ok := c.Transport.(*http.Transport)
		if !ok {
			return
		}
		t.Proxy = func(_ *http.Request) (*url.URL, error) {
			return proxy, nil
		}
	}
}

func WithProxyFunc(proxy func(*http.Request) (*url.URL, error)) Option {
	return func(c *Client) {
		t, ok := c.Transport.(*http.Transport)
		if !ok {
			return
		}
		t.Proxy = proxy
	}
}

func WithJar(jar http.CookieJar) Option {
	return func(c *Client) {
		c.Jar = jar
	}
}

func WithRoundTripper(t http.RoundTripper) Option {
	return func(c *Client) {
		c.Transport = t
	}
}

func WithDisableKeepAlives(disable bool) Option {
	return func(c *Client) {
		t, ok := c.Transport.(*http.Transport)
		if !ok {
			return
		}
		t.DisableKeepAlives = disable
	}
}

func WithTLSConfig(cfg *tls.Config) Option {
	return func(c *Client) {
		t, ok := c.Transport.(*http.Transport)
		if !ok {
			return
		}
		t.TLSClientConfig = cfg
	}
}

func WithTLSShuffle() Option {
	return WithTLSConfig(tlshuffle.NewConfig())
}

func WithRequestCallback(fn ...RequestCallback) Option {
	return func(c *Client) {
		c.requestCallbacks = append(c.requestCallbacks, fn...)
	}
}

func WithResponseCallback(fn ...ResponseCallback) Option {
	return func(c *Client) {
		c.responseCallbacks = append(c.responseCallbacks, fn...)
	}
}

func WithUserAgent(ua useragent.UserAgent) Option {
	str := useragent.String(ua)
	return WithRequestCallback(func(rq *http.Request) (*http.Request, error) {
		rq.Header.Set("User-Agent", str)
		return rq, nil
	})
}

func WithUserAgentString(ua string) Option {
	return WithRequestCallback(func(rq *http.Request) (*http.Request, error) {
		rq.Header.Set("User-Agent", ua)
		return rq, nil
	})
}

func WithCookies(url *url.URL, cookies []*http.Cookie) Option {
	return func(k *Client) {
		k.Jar.SetCookies(url, cookies)
	}
}

func WithLocator(l Locator) Option {
	return func(c *Client) {
		c.locator = l
	}
}
