package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr/conf"
	"github.com/hedzr/log"
	"github.com/hedzr/rate"
	"time"
)

// Config is an abstract structure for a rate limiter.
type Config struct {
	Name          string        `yaml:"name" json:"name,omitempty"`
	Description   string        `yaml:"description" json:"description,omitempty"`
	Enabled       bool          `yaml:"enabled" json:"enabled,omitempty"`
	Algorithm     string        `yaml:"algorithm" json:"algorithm,omitempty"`
	Interval      time.Duration `yaml:"interval" json:"interval,omitempty"`
	MaxRequests   int64         `yaml:"max-requests" json:"max-requests,omitempty"`
	HeaderKeyName string        `yaml:"header-key-name" json:"header-key-name,omitempty"`
	ExceptionKeys []string      `yaml:"exception-keys" json:"exception-keys,omitempty"`
	Routes        []string      `yaml:"routes" json:"routes,omitempty"` // Not Yet. Reserved for the future: [optional for routes builder]
}

// LoadConfig loads limiter config from cmdr config file and option store.
//
// keyPath is a dotted key-path string without:
// 1. cmdr Option Store Prefix (it's `app.` generally)
// 2. cmdr app-name Prefix (for most of config hierarchy)
//
// For example:
//
//    config := middleware.LoadConfig("server.rate-limite")
//    ginApp.Use(middleware.ForGin(config))
//
// And the corresponding config file should be:
//
// ```yaml
// app:
//   your-app: # <- replace it with your app name, the further KB in cmdr docs.
//     server:
//       rate-limits:
//         - name: by-api-key
//           interval: 1ms
//           max-requests: 30
//           header-key-name: X-API-KEY
//           exception-keys: [voxr-apps-test-api-key-fndsfjn]
// ```
//
// See also middleware.LoadConfigForGin
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

// LoadConfigForGin loads limiter config array from file (yaml/toml/...) and binds its into gin Router.
//
// A config file (eg. `rate-limit.yml`) should be put in cmdr-standard `conf.d` directory, so it can be loaded automatically:
//
// ```yaml
// app:
//   your-app: # <- replace it with your app name, the further KB in cmdr docs.
//     server:
//       rate-limits:
//         - name: by-api-key
//           interval: 1ms
//           max-requests: 30
//           header-key-name: X-API-KEY
//           exception-keys: [voxr-apps-test-api-key-fndsfjn]
// ```
//
func LoadConfigForGin(keyPath string, rg Router) {
	limiterConfigs := LoadConfig(keyPath)
	for _, cfg := range limiterConfigs {
		// if cfg.Enabled {
		rg.Use(ForGin(&cfg))
		// }
	}
}
