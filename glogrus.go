package glogrus

import (
	"fmt"
	"net/http"
	"time"
	"github.com/Sirupsen/logrus"
	"goji.io"
	"golang.org/x/net/context"
)

// NewGlogrus allows you to configure a goji middleware that logs all requests and responses
// using the structured logger logrus. It takes the logrus instance and the name of the app
// as the parameters and returns a middleware of type "func(goji.Handler) goji.Handler"
//
// Example:
//
//		package main
//
//		import(
//			""goji.io"
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
	return NewGlogrusWithReqId(l, name, emptyRequestId)
}

// NewGlogrusWithReqId allows you to configure a goji middleware that logs all requests and responses
// using the structured logger logrus. It takes the logrus instance, the name of the app and a function
// that can retrieve a requestId from the Context "func(context.Context) string"
// as the parameters and returns a middleware of type "func(goji.Handler) goji.Handler"
//
// Passing in the function that returns a requestId allows you to "plug in" other middleware that may set the request id
//
// Example:
//
//		package main
//
//		import(
//			""goji.io"
//			"github.com/goji/glogrus"
//			"github.com/Sirupsen/logrus"
//		)
//
//		func main() {
//
//			logr := logrus.New()
//			logr.Formatter = new(logrus.JSONFormatter)
//			goji.Use(glogrus.NewGlogrusWithReqId(logr, "my-app-name", GetRequestId))
//
//			goji.Get("/ping", yourHandler)
//			goji.Serve()
//		}
//
//		func GetRequestId(ctx context.Context) string {
//			return ctx.Value("requestIdKey")
//		}
//
func NewGlogrusWithReqId(l *logrus.Logger, name string, reqidf func(context.Context) string) func(goji.Handler) goji.Handler {
	return func(h goji.Handler) goji.Handler {
		fn := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			reqID := reqidf(ctx)

			l.WithFields(logrus.Fields{
				"req_id": reqID,
				"uri":    r.RequestURI,
				"method": r.Method,
				"remote": r.RemoteAddr,
			}).Info("req_start")
			lresp := wrapWriter(w)

			h.ServeHTTPC(ctx, lresp, r)
			lresp.maybeWriteHeader()

			latency := float64(time.Since(start)) / float64(time.Millisecond)

			l.WithFields(logrus.Fields{
				"req_id":  reqID,
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

func emptyRequestId(ctx context.Context) string {
	return ""
}
