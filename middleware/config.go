package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr/conf"
	"github.com/hedzr/log"
	"github.com/hedzr/rate"
	"time"
)

type Config struct {
	Name          string        `yaml:"name"`
	Description   string        `yaml:"description"`
	Algorithm     string        `yaml:"algorithm"`
	Interval      time.Duration `yaml:"interval"`
	MaxRequests   int64         `yaml:"max-requests"`
	HeaderKeyName string        `yaml:"header-key-name"`
	ExceptionKeys []string      `yaml:"exception-keys"`
	Routes        []string      `yaml:"routes"` // Not Yet. [optional for routes builder]
}

func LoadConfig(keyPath string) []Config {
	var dd []Config
	err := cmdr.GetSectionFrom(conf.AppName+"."+keyPath, &dd)
	if err != nil {
		log.Warnf("load '%v' failed: %v", keyPath, err)
	}
	for _, c := range dd {
		if c.Algorithm == "" {
			c.Algorithm = string(rate.TokenBucket)
		}
	}
	return dd
}

type Router interface {
	gin.IRouter
	// gin.IRoutes
}

func LoadConfigForGin(keyPath string, rg Router) {
	limiterConfigs := LoadConfig(keyPath)
	for _, cfg := range limiterConfigs {
		rg.Use(ForGin(&cfg))
	}
}
