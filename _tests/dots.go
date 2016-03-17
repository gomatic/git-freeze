package main

import (
	"crypto/x509/pkix"
	"net/http/cookiejar"
	"net/http/pprof"

	"github.com/Spatially/go-flagged"
)

func main() {
	flagged.Parse(nil)
	pkix.Name
	cookiejar.Options
	pprof.Profile(nil, nil)
}
