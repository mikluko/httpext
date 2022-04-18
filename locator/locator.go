package locator

import "net"

type Location struct {
	IP      net.IP
	ASN     string
	Country string
	Region  string
	City    string
	Postal  string
}
