// basic frontend
// Under MIT license see LICENSE file.
package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"

	"github.com/golang/glog"
	"github.com/vulcand/oxy/forward"
	"github.com/vulcand/oxy/roundrobin"
	"github.com/vulcand/oxy/stream"
)

type stringslice []string

func (b *stringslice) String() string {
	return fmt.Sprintf("%v", *b)
}

func (b *stringslice) Set(value string) error {
	*b = append(*b, value)
	return nil
}

var backend stringslice

type logger struct {
}

func (l logger) Infof(format string, args ...interface{}) {
	glog.Infof(format, args)
}
func (l logger) Warningf(format string, args ...interface{}) {
	glog.Warningf(format, args)
}

func (l logger) Errorf(format string, args ...interface{}) {
	glog.Errorf(format, args)
}

func init() {
	flag.Var(&backend, "backend", "backend url")
}

func main() {

	flag.Parse()
	if len(backend) == 0 {
		flag.PrintDefaults()
		return
	}
	glog.Infof("Starting - using %v", backend)

	// Forwards incoming requests to whatever location URL points to, adds proper forwarding headers
	fwd, err := forward.New()
	if err != nil {
		glog.Fatal(err)
	}

	lb, err := roundrobin.New(fwd)
	if err != nil {
		glog.Fatal(err)
	}

	stream, err := stream.New(lb, stream.Retry(`IsNetworkError() && Attempts() < 4`), stream.Logger(logger{}))
	if err != nil {
		glog.Fatal(err)
	}

	for _, b := range backend {
		u, err := url.Parse(b)
		if err != nil {
			glog.Fatal(err)
		}

		lb.UpsertServer(u)
	}

	// that's it! our reverse proxy is ready!
	s := &http.Server{
		Addr:    ":8080",
		Handler: stream,
	}
	err = s.ListenAndServe()
	if err != nil {
		glog.Fatal(err)
	}

}
