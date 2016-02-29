package glogrus

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"goji.io"
	"golang.org/x/net/context"
)

// glogrus is a middleware handler that logs the
// request and response in a structured way
type glogrus struct {
	h    goji.Handler
	c    context.Context
	l    *logrus.Logger
	name string
}

func (glogr glogrus) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	start := time.Now()
	
	//TODO: figure out how to get a proper reqId in the context
	//reqID := middleware.GetReqID(*glogr.c)
	
	glogr.l.WithFields(logrus.Fields{
		//"req_id": reqID,
		"uri":    req.RequestURI,
		"method": req.Method,
		"remote": req.RemoteAddr,
	}).Info("req_start")
	lresp := wrapWriter(resp)

	glogr.h.ServeHTTPC(glogr.c, lresp, req)
	lresp.maybeWriteHeader()

	latency := float64(time.Since(start)) / float64(time.Millisecond)

	glogr.l.WithFields(logrus.Fields{
		//"req_id":  reqID,
		"status":  lresp.status(),
		"method":  req.Method,
		"uri":     req.RequestURI,
		"remote":  req.RemoteAddr,
		"latency": fmt.Sprintf("%6.4f ms", latency),
		"app":     glogr.name,
	}).Info("req_served")
}

// NewGlogrus allows you to configure a goji middleware that logs all requests and responses
// using the structured logger logrus. It takes the logrus instance and the name of the app
// as the parameters and returns a middleware of type "func(c *web.C, http.Handler) http.Handler"
//
// Example:
//
//		package main
//
//		import(
//			"github.com/zenazn/goji"
//			"github.com/zenazn/goji/web/middleware"
//			"github.com/goji/glogrus"
//			"github.com/Sirupsen/logrus"
//		)
//
//		func main() {
//
//			logr := logrus.New()
//			logr.Formatter = new(logrus.JSONFormatter)
//			goji.Use(glogrus.NewGlogrus(logr, "my-app-name"))
//
//			goji.Get("/ping", yourHandler)
//			goji.Serve()
//		}
//
func NewGlogrus(l *logrus.Logger, name string) func(goji.Handler) goji.Handler {
	return func(h goji.Handler) goji.Handler {
		fn := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			//TODO: figure out how to get a proper reqId in the context
			//reqID := middleware.GetReqID(*glogr.c)

			l.WithFields(logrus.Fields{
				//"req_id": reqID,
				"uri":    r.RequestURI,
				"method": r.Method,
				"remote": r.RemoteAddr,
			}).Info("req_start")
			lresp := wrapWriter(w)

			h.ServeHTTPC(ctx, lresp, r)
			lresp.maybeWriteHeader()

			latency := float64(time.Since(start)) / float64(time.Millisecond)

			l.WithFields(logrus.Fields{
				//"req_id":  reqID,
				"status":  lresp.status(),
				"method":  r.Method,
				"uri":     r.RequestURI,
				"remote":  r.RemoteAddr,
				"latency": fmt.Sprintf("%6.4f ms", latency),
				"app":     name,
			}).Info("req_served")
		}
		return goji.HandlerFunc(fn)
	}

}
