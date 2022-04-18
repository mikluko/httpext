package httpext

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"

	"github.com/mikluko/httpext/locator"
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

type testLocator struct {
	mock.Mock
}

func (t *testLocator) Locate(ctx context.Context, c http.RoundTripper) (*locator.Location, error) {
	args := t.Called(ctx, c)
	x := args.Get(0)
	if x != nil {
		return x.(*locator.Location), nil
	}
	return nil, args.Error(1)
}

func TestClient_Location(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	l := testLocator{}
	c := New(WithLocator(&l))

	loc0 := &locator.Location{}

	l.On("Locate", ctx, c.Transport).Return(loc0, nil).Once()

	loc1, err := c.Location(ctx)
	require.NoError(t, err)
	require.Same(t, loc0, loc1)
}

func TestClient_LocationCached(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	l := testLocator{}
	c := New(WithLocator(&l))

	loc0 := &locator.Location{}

	l.On("Locate", ctx, c.Transport).Return(loc0, nil).Once()

	loc1, err := c.LocationCached(ctx)
	require.NoError(t, err)
	require.Same(t, loc0, loc1)

	loc2, err := c.LocationCached(ctx)
	require.NoError(t, err)
	require.Same(t, loc0, loc2)
}
