# goji/glogrus [![GoDoc](https://godoc.org/github.com/goji/glogrus?status.png)](https://godoc.org/github.com/goji/glogrus) [![Build Status](https://travis-ci.org/goji/glogrus.svg)](https://travis-ci.org/goji/glogrus)

glogrus provides structured logging via logrus for Goji. 

## Example


```go

package main

import (
	"github.com/goji/glogrus"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web/middleware"
)

func main() {
	goji.Abandon(middleware.Logger)
	goji.Use(glogrus.New())

	goji.Get("/ping", yourHandler)
	goji.Serve()
}

```

Consider using glogrus.Custom() for more control over the logrus instance and service name


#### Looking for hierarchical structured logging?
[slog](https://github.com/zenazn/slog) and [lunk](https://github.com/codahale/lunk) looks interesting.