package httpext

import (
	"errors"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"

	"github.com/hashicorp/go-cleanhttp"
	"golang.org/x/net/context"
	"golang.org/x/net/publicsuffix"

	"github.com/mikluko/httpext/locator"
)

func New(opts ...Option) *Client {
	c := Client{
		Client: cleanhttp.DefaultClient(),
	}
	c.Jar = defaultJar()
	c.locator = defaultLocator()
	for i := range opts {
		opts[i](&c)
	}
	return &c
}

type Client struct {
	*http.Client
	requestCallbacks  []RequestCallback
	responseCallbacks []ResponseCallback
	locator           Locator
	location          *locator.Location
	locationMux       sync.Mutex
}

func (c *Client) Do(rq *http.Request) (*http.Response, error) {
	var err error
	for i := range c.requestCallbacks {
		rq, err = c.requestCallbacks[i](rq)
		if err != nil {
			return nil, err
		}
	}
	var rs *http.Response
	rs, err = c.Client.Do(rq)
	if err != nil {
		return nil, err
	}
	for i := range c.responseCallbacks {
		rs, err = c.responseCallbacks[i](rs)
		if err != nil {
			return nil, err
		}
	}
	return rs, nil
}

func (c *Client) Get(url string) (*http.Response, error) {
	rq, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(rq)
}

func (c *Client) Head(url string) (*http.Response, error) {
	rq, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(rq)
}

func (c *Client) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	rq, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	rq.Header.Set("Content-Type", contentType)
	return c.Do(rq)
}

func (c *Client) PostForm(url string, data url.Values) (*http.Response, error) {
	return c.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

var ErrCookieNotFound = errors.New("cookie not found")

func (c *Client) Cookie(url *url.URL, name string) (*http.Cookie, error) {
	cookies := c.Jar.Cookies(url)
	for i := range cookies {
		if cookies[i].Name == name {
			return cookies[i], nil
		}
	}
	return nil, ErrCookieNotFound
}

func (c *Client) HasCookie(url *url.URL, name string) bool {
	cookie, err := c.Cookie(url, name)
	if err != nil {
		return false
	}
	if cookie == nil {
		panic("programming error")
	}
	return true
}

func (c *Client) CookieValue(url *url.URL, name string) (string, bool) {
	cookie, err := c.Cookie(url, name)
	if err != nil {
		return "", false
	}
	if cookie == nil {
		panic("programming error")
	}
	return cookie.Value, true
}

func (c *Client) Location(ctx context.Context) (*locator.Location, error) {
	return c.locator.Locate(ctx, c.Transport)
}

func (c *Client) LocationCached(ctx context.Context) (*locator.Location, error) {
	c.locationMux.Lock()
	defer c.locationMux.Unlock()
	if c.location != nil {
		return c.location, nil
	}
	loc, err := c.Location(ctx)
	if err != nil {
		return nil, err
	}
	c.location = loc
	return loc, nil
}

type RequestCallback func(rq *http.Request) (*http.Request, error)

type ResponseCallback func(rs *http.Response) (*http.Response, error)

type Locator interface {
	Locate(context.Context, http.RoundTripper) (*locator.Location, error)
}

func defaultLocator() Locator {
	return &locator.Ipinfo{}
}

func defaultJar() *cookiejar.Jar {
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	return jar
}