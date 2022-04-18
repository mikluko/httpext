package httpext

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	c := New()
	assert.NotNil(t, c)
	assert.NotNil(t, c.Jar)
	assert.NotNil(t, c.Transport)
}

func TestClient_HasCookie(t *testing.T) {
	c := New()
	u, _ := url.Parse("http://example.com/")
	c.Jar.SetCookies(u, []*http.Cookie{
		{Name: "example_name", Value: "example_value"},
	})

	t.Run("set", func(t *testing.T) {
		assert.True(t, c.HasCookie(u, "example_name"))
	})
	t.Run("unset", func(t *testing.T) {
		assert.False(t, c.HasCookie(u, "other_name"))
	})
}

func TestClient_CookieValue(t *testing.T) {
	c := New()
	u, _ := url.Parse("http://example.com/")
	c.Jar.SetCookies(u, []*http.Cookie{
		{Name: "example_name", Value: "example_value"},
	})

	t.Run("set", func(t *testing.T) {
		value, isSet := c.CookieValue(u, "example_name")
		assert.True(t, isSet)
		assert.Equal(t, "example_value", value)
	})
	t.Run("unset", func(t *testing.T) {
		value, isSet := c.CookieValue(u, "other_name")
		assert.False(t, isSet)
		assert.Equal(t, "", value)
	})
}
