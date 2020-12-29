package limiter

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/AyushSenapati/limiter"

	"github.com/AyushSenapati/guardian/lib/proxy"
)

const (
	limiterHeaderLimit     = "X-RateLimit-Limit"
	limiterHeaderRemaining = "X-RateLimit-Remaining"
)

// Config defines limiter config
type Config struct {
	Quota int64
	Per   string
	Store string
}

// SetupLimiter implements the logic to read the provided raw config and configure itself
func SetupLimiter(def *proxy.RouterDefinition, rawConfig map[string]interface{}) error {
	// var config Config

	// validJSON, err := json.Marshal(rawConfig)
	// if err != nil {
	// 	return err
	// }
	// if err := json.Unmarshal(validJSON, &config); err != nil {
	// 	return err
	// }
	config, err := limiter.ParseAndLoadConfiguration(rawConfig)
	if err != nil {
		return err
	}
	l := limiter.NewWithConfiguration(config)
	l.Configure()
	def.AddMiddleware(l.Handler())
	// def.AddMiddleware(configureLimiterMW(config))

	return nil
}

func configureLimiterMW(config Config) func(http.Handler) http.Handler {
	store := NewMemoryStore(config)
	storeCleaner.addStore(store)
	storeCleaner.dispatchCleaner()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get source IP from request header
			ipAddr := getIP(r)
			if ipAddr == "" {
				w.Write([]byte("Malformed request IP detected"))
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			w.Header().Set(
				limiterHeaderLimit, fmt.Sprintf("%d/%s", config.Quota, config.Per),
			)
			quotaRemaining, isAllowed := store.IsAllowed(ipAddr)
			w.Header().Set(limiterHeaderRemaining, strconv.Itoa(int(quotaRemaining)))
			if !isAllowed {
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
