# go-rate

![Go](https://github.com/hedzr/rate/workflows/Go/badge.svg)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/hedzr/rate.svg?label=release)](https://github.com/hedzr/rate/releases)
[![](https://img.shields.io/badge/go-dev-green)](https://pkg.go.dev/github.com/hedzr/rate)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/hedzr/rate)
<!-- [![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fhedzr%2Frate.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fhedzr%2Frate?ref=badge_shield) -->
[![Go Report Card](https://goreportcard.com/badge/github.com/hedzr/rate)](https://goreportcard.com/report/github.com/hedzr/rate)
[![Coverage Status](https://coveralls.io/repos/github/hedzr/rate/badge.svg?branch=master&.9)](https://coveralls.io/github/hedzr/rate?branch=master)

`go-rate` provides the rate limiters generally.

## Usages

### Simple

```go
package main

import (
	"fmt"
	"github.com/hedzr/rate"
	"time"
)

func main() {
	l := rate.New(rate.LeakyBucket, 100, time.Second)
	for i := 0; i < 120; i++ {
		ok := l.Take(1)
		if !ok {
			fmt.Printf("#%d Take() returns not ok, counter: %v\n", i, rate.CountOf(l))
			time.Sleep(50 * time.Millisecond)
		}
	}
}
```

### As a gin middleware

```go
package main

import (
   "github.com/gin-gonic/gin"
   "github.com/hedzr/rate"
   "github.com/hedzr/rate/middleware"
   "time"
)

func webserver() {
	r := engine()
	r.Run(":3000")
}

func engine() *gin.Engine {
	config := &middleware.Config{
		Name:          "...",
		Description:   "...",
		Algorithm:     string(rate.TokenBucket),
		Interval:      time.Second,
		MaxRequests:   1000,
		HeaderKeyName: "X-API-TOKEN",
		ExceptionKeys: nil,
		Routes:        nil,
	}
	r := gin.Default()
	r.Use(middleware.ForGin(config))
	return r
}
```

### Load limit config with cmdr Option Store

While integrated with [hedzr/cmdr](https://github.com/hedzr/cmdr), the short loading is available:

```go
import "github.com/hedzr/rate/middleware"

func BuildRoutes(rg *gin.Engine) *gin.Engine {
    buildRoutes(rg.Group("/prefix"), rg)
}
func buildRoutes(rg Router, root *gin.Engine) {
    middleware.LoadConfigForGin("server.rate-limits", rg)
    rg.Get("/echo/*action", echoHandler)
}
func echoGinHandler(c *gin.Context) {
    action := c.Param("action")
    if action == "" || action == "/" {
        action = "<no action>"
    }
    _, _ = io.WriteString(c.Writer, fmt.Sprintf("action: %v\n", action))
}
```

A config file (eg. `rate-limit.yml`) should be put in cmdr-standard `conf.d` directory, so it can be loaded automatically:

```yaml
app:
  your-app: # <- replace it with your app name, the further KB in cmdr docs.
    server:
      rate-limits:
        - name: by-api-key
          interval: 1ms
          max-requests: 30
          header-key-name: X-API-KEY
          exception-keys: [voxr-apps-test-api-key-fndsfjn]
```



## License

MIT

