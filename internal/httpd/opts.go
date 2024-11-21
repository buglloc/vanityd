package httpd

import "strings"

const (
	DefaultAddr = "localhost:3000"
	DefaultZone = "prj.buglloc.com"
)

type Option func(h *HttpD)

func WithAddr(addr string) Option {
	return func(h *HttpD) {
		h.http.Addr = addr
	}
}

func WithZone(zone string) Option {
	return func(h *HttpD) {
		h.zone = "." + strings.Trim(zone, ". \t")
	}
}
