// Copyright 2013 Nikita Kovaliov, maizy.ru
// See MIT-LICENSE.txt for details.

package errserv

import (
	"net/http"
)

type Port uint16
type Errcode uint16
type Timeout uint16

//port =>error_code
type ErrorPortPair map[Port]Errcode

//port =>timeout in ms
type TimeoutPortPair map[Port]Timeout

type Options struct {
	ErrorsPorts  ErrorPortPair
	TimeoutPorts TimeoutPortPair
}

func (opt *Options) IsEnableErrorServ() bool {
	return len(opt.ErrorsPorts) > 0
}

func (opt *Options) IsEnableTimeoutServ() bool {
	return len(opt.TimeoutPorts) > 0
}

func NewOptions() *Options {
	return &Options{
		ErrorsPorts:  make(ErrorPortPair),
		TimeoutPorts: make(TimeoutPortPair),
	}
}

func parseFlags() *Options {
	// TODO: parse config file or args
	opt := NewOptions()
	var baseErrorPort uint16 = 10000
	stdErrors := []uint16{
		http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden,
		http.StatusNotFound, http.StatusMethodNotAllowed, http.StatusNotAcceptable,
		http.StatusGone, http.StatusLengthRequired, http.StatusPreconditionFailed,
		http.StatusRequestEntityTooLarge, http.StatusRequestURITooLong, http.StatusUnsupportedMediaType,
		http.StatusRequestedRangeNotSatisfiable,

		http.StatusInternalServerError, http.StatusNotImplemented, http.StatusBadGateway,
		http.StatusServiceUnavailable, http.StatusGatewayTimeout,
	}

	for _, errcode := range stdErrors {
		opt.ErrorsPorts[Port(baseErrorPort+uint16(errcode))] = Errcode(errcode)
	}
	opt.TimeoutPorts[11010] = 1000
	opt.TimeoutPorts[11012] = 1200
	opt.TimeoutPorts[11015] = 1500
	opt.TimeoutPorts[11020] = 2000
	opt.TimeoutPorts[11050] = 5000
	opt.TimeoutPorts[11080] = 8000
	opt.TimeoutPorts[11100] = 10000
	return opt
}
