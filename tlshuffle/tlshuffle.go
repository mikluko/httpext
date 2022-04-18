package tlshuffle

import (
	"crypto/tls"
	"math/rand"
	"time"
)

var cipherSuitesMain = []uint16{
	tls.TLS_AES_128_GCM_SHA256,
	tls.TLS_AES_256_GCM_SHA384,
	tls.TLS_CHACHA20_POLY1305_SHA256,
}

var cipherSuitesExtra = []uint16{
	tls.TLS_RSA_WITH_RC4_128_SHA,
	tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
	tls.TLS_RSA_WITH_AES_128_CBC_SHA,
	tls.TLS_RSA_WITH_AES_256_CBC_SHA,
	tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
	tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
	tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
	tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
	tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
	tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
	tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
	tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
	tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
	tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
	tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
}

var defaultSuites = append(cipherSuitesMain, cipherSuitesExtra...)

func DefaultSuites() []uint16 {
	suites := make([]uint16, len(defaultSuites))
	copy(suites, defaultSuites)
	return suites
}

func init() {
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
}

var rnd *rand.Rand

func ShuffleCipherSuites(suites []uint16) {
	if len(suites) < 3 {
		return
	}
	if rnd.Intn(2) == 1 {
		suites[1], suites[2] = suites[2], suites[1]
	}
	if len(suites) < 5 {
		return
	}
	rnd.Shuffle(len(suites[3:]), func(i, j int) {
		suites[i+3], suites[j+3] = suites[j+3], suites[i+3]
	})
}

func CipherSuites() []uint16 {
	suites := make([]uint16, len(defaultSuites))
	copy(suites, defaultSuites)
	ShuffleCipherSuites(suites)
	return suites
}

func ShuffleConfig(cfg *tls.Config) {
	ShuffleCipherSuites(cfg.CipherSuites)
}

func NewConfig() *tls.Config {
	return &tls.Config{
		CipherSuites: CipherSuites(),
	}
}
