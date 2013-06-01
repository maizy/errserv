// Copyright 2013 Nikita Kovaliov, maizy.ru
// See MIT-LICENSE.txt for details.

package errserv

import (
	"fmt"
	"log"
	"net/http"
)

const (
	bindFailed  = 0
	bindSuccess = 1
)

type PortHandler interface {
	Port() Port
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type OneErrorHandler struct {
	port  Port
	error Errcode
}

func (handler *OneErrorHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	error := int(handler.error)
	log.Printf("[REQUEST :%d] %s %s => %d", handler.port, request.Method, request.RequestURI, handler.error)
	http.Error(writer, fmt.Sprintf("Error %d - %s", error, http.StatusText(error)), error)
}

func (handler *OneErrorHandler) Port() Port {
	return handler.port
}

func Main() {
	opt := parseFlags()

	bindChan := make(chan int)
	bind := func(handler PortHandler) {
		port := handler.Port()
		err := http.ListenAndServe(fmt.Sprintf(":%d", port), http.Handler(handler))
		if err != nil {
			log.Printf("[WARN] Port %d skipping: %s", port, err)
			bindChan <- bindFailed
		} else {
			bindChan <- bindSuccess
		}
	}

	if !opt.IsEnableErrorServ() && !opt.IsEnableTimeoutServ() {
		log.Fatal("Nothing to do")
	}

	needBinded := 0
	if opt.IsEnableErrorServ() {
		for port, errcode := range opt.ErrorsPorts {
			needBinded++
			go bind(&OneErrorHandler{error: errcode, port: port})
		}
	}

	anyBind := false
	var res int
	for i := 0; i < needBinded; i++ {
		res = <-bindChan
		anyBind = anyBind || res == bindSuccess
	}
	if !anyBind {
		log.Fatal("All ports are skipped")
	}
}
