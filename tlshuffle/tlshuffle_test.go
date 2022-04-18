package tlshuffle

import (
	"crypto/tls"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultSuites(t *testing.T) {
	t.Run("clone works", func(t *testing.T) {
		suites := DefaultSuites()
		assert.Len(t, suites, len(defaultSuites))
		for i := range suites {
			assert.Equal(t, suites[i], defaultSuites[i])
		}
	})
	t.Run("source is immutable", func(t *testing.T) {
		suites := DefaultSuites()
		suites[0], suites[1] = suites[1], suites[0]
		require.Equal(t, defaultSuites[1], suites[0])
		require.Equal(t, defaultSuites[0], suites[1])
	})
}

func testSuites(t *testing.T, orig, shuf []uint16) {
	require.Len(t, shuf, len(orig))
	require.GreaterOrEqual(t, len(orig), 3)
	assert.Equal(t, orig[0], shuf[0])
	assert.Contains(t, orig[1:3], shuf[1])
	assert.Contains(t, orig[1:3], shuf[2])
	for i := range shuf {
		if shuf[i] != orig[i] {
			return
		}
	}
	assert.Fail(t, "suites don't seem to be shuffled")
}

func TestShuffleCipherSuites(t *testing.T) {
	suites := DefaultSuites()
	ShuffleCipherSuites(suites)
	testSuites(t, defaultSuites, suites)
}

func TestCipherSuites(t *testing.T) {
	suites := CipherSuites()
	testSuites(t, defaultSuites, suites)
}

func TestNewConfig(t *testing.T) {
	cfg := NewConfig()
	testSuites(t, defaultSuites, cfg.CipherSuites)
}

func TestShuffleConfig(t *testing.T) {
	cfg := tls.Config{
		CipherSuites: CipherSuites(),
	}
	ShuffleConfig(&cfg)
	testSuites(t, defaultSuites, cfg.CipherSuites)
}
