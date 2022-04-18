package locator

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"

	"golang.org/x/net/context"
)

const ipinfoEndpoint = "http://ipinfo.io"

type Ipinfo struct {
	Token string
}

func (l *Ipinfo) Locate(ctx context.Context, c http.RoundTripper) (*Location, error) {
	if c == nil {
		c = http.DefaultTransport
	}
	rq, _ := http.NewRequestWithContext(ctx, http.MethodGet, ipinfoEndpoint, nil)
	if l.Token != "" {
		rq.Header.Add("Authorization", "Bearer "+l.Token)
	}
	rs, err := c.RoundTrip(rq)
	if err != nil {
		return nil, err
	}
	defer rs.Body.Close()

	var data ipinfoData
	err = json.NewDecoder(rs.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data.ToLocation(), nil
}

type ipinfoData struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Loc      string `json:"loc"`
	Org      string `json:"org"`
	Postal   string `json:"postal"`
	Timezone string `json:"timezone"`
}

func (d *ipinfoData) asn() string {
	bits := strings.SplitN(d.Org, " ", 2)
	return bits[0]
}

func (d *ipinfoData) ToLocation() *Location {
	return &Location{
		IP:      net.ParseIP(d.IP),
		ASN:     d.asn(),
		Country: d.Country,
		Region:  d.Region,
		City:    d.City,
		Postal:  d.Postal,
	}
}
