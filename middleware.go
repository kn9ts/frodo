package frodo

import (
	"net/http"
	"strings"
)

// Request will help facilitate the passing of multiple handlers
type Request struct {
	handlers       []Handle
	total          int
	nextPosition   int
	ResponseWriter *AppResponseWriter
	*http.Request
	Params
}

// Middleware declares the minimum implementation necessary for a handlers
// to be used as Frodo's middleware route Handlers
type Middleware interface {
	Next(w *AppResponseWriter)
}

func (r *Request) runHandleChain(w *AppResponseWriter) {
	r.nextPosition++
	r.ResponseWriter = w
	r.handlers[0](w, r)
}

// Next will be used to call the next handler in line/queue
func (r *Request) Next() {
	// 1st check if the next handler position accounts for the number
	// of handlers existing in the handlers array
	if r.nextPosition < r.total {
		// get the next handler
		nextHandler := r.handlers[r.nextPosition]
		// move the cursor
		r.nextPosition++

		// 1st check if a write has happened
		// meaning a response has been issued out to the client
		// if not run the next handler in line
		if r.ResponseWriter.WriteHappened() == false {
			nextHandler(r.ResponseWriter, r)
		}
	}
	return
}

// ClientIP implements a best effort algorithm to return the real client IP, it parses
// X-Real-IP and X-Forwarded-For in order to work properly with reverse-proxies such us: nginx or haproxy.
func (r *Request) ClientIP() string {
	if true {
		clientIP := strings.TrimSpace(r.Request.Header.Get("X-Real-Ip"))
		if len(clientIP) > 0 {
			return clientIP
		}
		clientIP = r.Request.Header.Get("X-Forwarded-For")
		if index := strings.IndexByte(clientIP, ','); index >= 0 {
			clientIP = clientIP[0:index]
		}
		clientIP = strings.TrimSpace(clientIP)
		if len(clientIP) > 0 {
			return clientIP
		}
	}
	return strings.TrimSpace(r.Request.RemoteAddr)
}
