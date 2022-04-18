package httpext

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mikluko/httpext/useragent"
)

type cookieJarMock struct{}

func (c cookieJarMock) SetCookies(_ *url.URL, _ []*http.Cookie) {
	panic("implement me")
}

func (c cookieJarMock) Cookies(u *url.URL) []*http.Cookie {
	panic("implement me")
}

func TestWithJar(t *testing.T) {
	jar := &cookieJarMock{}
	c := New(WithJar(jar))
	require.Same(t, jar, c.Jar.(*cookieJarMock))
}

type roundTripperMock struct{}

func (r roundTripperMock) RoundTrip(_ *http.Request) (*http.Response, error) {
	panic("implement me")
}

func TestWithRoundTripper(t *testing.T) {
	r := &roundTripperMock{}
	c := New(WithRoundTripper(r))
	require.Same(t, r, c.Transport.(*roundTripperMock))
}

func TestWithProxyURL(t *testing.T) {
	p0, _ := url.Parse("http://example.com:8080")
	c := New(WithProxyURL(p0))

	rq, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	p1, err := c.Transport.(*http.Transport).Proxy(rq)
	require.NoError(t, err)
	require.Same(t, p0, p1)
}

func TestWithProxyFunc(t *testing.T) {
	p0, _ := url.Parse("http://example.com:8080")
	f := func(_ *http.Request) (*url.URL, error) {
		return p0, nil
	}
	c := New(WithProxyFunc(f))

	rq, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	p1, err := c.Transport.(*http.Transport).Proxy(rq)
	require.NoError(t, err)
	require.Same(t, p0, p1)
}

func TestWithTLSConfig(t *testing.T) {
	cfg := &tls.Config{}
	c := New(WithTLSConfig(cfg))
	require.Same(t, cfg, c.Transport.(*http.Transport).TLSClientConfig)
}

func TestWithTLSShuffle(t *testing.T) {
	cfg := &tls.Config{}
	c := New(WithTLSConfig(cfg), WithTLSShuffle()) // it should replace earlier one
	require.NotSame(t, cfg, c.Transport.(*http.Transport).TLSClientConfig)
}

func TestWithDisableKeepAlives(t *testing.T) {
	c := New(WithDisableKeepAlives(true))
	require.True(t, c.Transport.(*http.Transport).DisableKeepAlives)
}

func TestWithRequestCallback(t *testing.T) {
	rq0, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	c := New(WithRequestCallback(func(_ *http.Request) (*http.Request, error) {
		return rq0, nil
	}))
	assert.Len(t, c.requestCallbacks, 1)
	rq1, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	rq2, err := c.requestCallbacks[0](rq1)
	require.NoError(t, err)
	require.Same(t, rq0, rq2)
}

func TestWithResponseCallback(t *testing.T) {
	rs0 := &http.Response{}
	c := New(WithResponseCallback(func(_ *http.Response) (*http.Response, error) {
		return rs0, nil
	}))
	assert.Len(t, c.responseCallbacks, 1)
	rs1 := &http.Response{}
	rs2, err := c.responseCallbacks[0](rs1)
	require.NoError(t, err)
	require.Same(t, rs0, rs2)
}

func TestWithUserAgent(t *testing.T) {
	c := New(WithUserAgent(useragent.Any))
	assert.Len(t, c.requestCallbacks, 1)
	rq0, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	rq1, err := c.requestCallbacks[0](rq0)
	require.NoError(t, err)
	require.NotEqual(t, "", rq1.Header.Get("User-Agent"))
}

func TestWithUserAgentString(t *testing.T) {
	c := New(WithUserAgentString("bla-bla"))
	assert.Len(t, c.requestCallbacks, 1)
	rq0, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	rq1, err := c.requestCallbacks[0](rq0)
	require.NoError(t, err)
	require.Equal(t, "bla-bla", rq1.Header.Get("User-Agent"))
}

func TestWithCookies(t *testing.T) {
	u, _ := url.Parse("http://example.com/")
	c := New(WithCookies(u,
		[]*http.Cookie{
			{Name: "test_name", Value: "test_value"},
		}),
	)
	cookies := c.Jar.Cookies(u)
	require.Len(t, cookies, 1)
	require.Equal(t, "test_name", cookies[0].Name)
	require.Equal(t, "test_value", cookies[0].Value)
}

func TestWithLocator(t *testing.T) {
	l := testLocator{}
	c := New(WithLocator(&l))
	require.Same(t, c.locator, &l)
}
