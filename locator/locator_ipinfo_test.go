//go:build integration

package locator

import (
	"context"
	"fmt"
	"os"
)

func ExampleIpinfo_Location() {
	l := Ipinfo{
		Token: os.Getenv("IPINFO_TOKEN"),
	}
	loc, err := l.Locate(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(loc)
	// Output:
}
