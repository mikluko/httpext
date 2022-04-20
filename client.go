package httpext

import (
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"

	"github.com/hashicorp/go-cleanhttp"
	"golang.org/x/net/publicsuffix"
)

func New(opts ...Option) *Client {
	c := Client{
		Client: cleanhttp.DefaultClient(),
	}
	c.Jar = defaultJar()
	for i := range opts {
		opts[i](&c)
	}
	return &c
}

type Client struct {
	*http.Client
	requestCallbacks  []RequestCallback
	responseCallbacks []ResponseCallback
	ownIP             net.IP
	ownIPMux          sync.Mutex
}

func (c *Client) Do(rq *http.Request) (*http.Response, error) {
	var err error
	for i := range c.requestCallbacks {
		rq, err = c.requestCallbacks[i](c, rq)
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
		rs, err = c.responseCallbacks[i](c, rs)
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

type RequestCallback func(c *Client, rq *http.Request) (*http.Request, error)

type ResponseCallback func(c *Client, rs *http.Response) (*http.Response, error)

func defaultJar() *cookiejar.Jar {
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	return jar
}
