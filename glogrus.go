package glogrus

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

// glogrus is a middleware handler that logs the
// request and response in a structured way
type glogrus struct {
	h    http.Handler
	c    *web.C
	l    *logrus.Logger
	name string
}

func (glogr glogrus) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	start := time.Now()
	reqID := middleware.GetReqID(*glogr.c)
	glogr.l.WithFields(logrus.Fields{
		"req_id": reqID,
		"uri":    req.RequestURI,
		"method": req.Method,
		"remote": req.RemoteAddr,
	}).Info("req_start")
	lresp := wrapWriter(resp)

	glogr.h.ServeHTTP(lresp, req)
	lresp.maybeWriteHeader()

	latency := float64(time.Since(start)) / float64(time.Millisecond)

	glogr.l.WithFields(logrus.Fields{
		"req_id":  reqID,
		"status":  lresp.status(),
		"method":  req.Method,
		"uri":     req.RequestURI,
		"remote":  req.RemoteAddr,
		"latency": fmt.Sprintf("%7.7f ms", latency),
		"app":     glogr.name,
	}).Info("req_served")
}

// Custom takes a logrus instace and name of the service
// as parameters and returns glogrus http.Handler
func Custom(l *logrus.Logger, name string) func(*web.C, http.Handler) http.Handler {
	fn := func(c *web.C, h http.Handler) http.Handler {
		return glogrus{h: h, c: c, l: l, name: name}
	}
	return fn
}

// New returns a glogrus http.Handler
func New() func(*web.C, http.Handler) http.Handler {
	fn := func(c *web.C, h http.Handler) http.Handler {
		log := logrus.New()
		log.Level = logrus.Info
		log.Formatter = &logrus.TextFormatter{}
		return glogrus{h: h, c: c, l: log, name: "*"}
	}
	return fn
}
