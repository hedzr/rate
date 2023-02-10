//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hedzr/rate"
	"github.com/hedzr/rate/supports/middleware"
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
