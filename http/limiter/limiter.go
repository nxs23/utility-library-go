package limiter

import (
	"net/http"
	"os"
	"strconv"
	"time"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
	"github.com/nxs23/utility-library-go/cache"
	log "github.com/nxs23/utility-library-go/logging/logger"
)

var logger *log.Logger

type limiter struct {
	Store ratelimit.Store
}

func init() {
	// Create new Logger
	logger = log.NewLogger("limiter")
}

func newLimiter(config Config) *limiter {
	cache := cache.Default()
	// Ping cache to confirm that is working
	cache.Ping()
	store := newStore(config, cache)

	return &limiter{
		Store: store,
	}
}

func (limiter *limiter) ApplyRateLimiter() gin.HandlerFunc {

	rm := ratelimit.RateLimiter(limiter.Store, &ratelimit.Options{
		ErrorHandler: errorHandler,
		KeyFunc:      keyFunc,
	})

	return rm
}

func Default() gin.HandlerFunc {
	// Check if rate limit has been reconfigured by ENV
	limit := RATE_LIMITER_LIMIT
	l := os.Getenv("RATE_LIMITER_LIMIT")

	// Override default config with ENV configuration
	if l != "" {
		logger.Printf("Read rate limit from RATE_LIMITER_LIMIT env: %s", l)
		limit, _ = strconv.Atoi(l)
	}

	// Check if rate timeframe has been reconfigured by ENV
	timeframe := RATE_LIMITER_TIMEFRAME
	t := os.Getenv("RATE_LIMITER_TIMEFRAME")

	// Override default config with ENV configuration
	if t != "" {
		logger.Printf("Read rate timeframe from RATE_LIMITER_TIMEFRAME env: %s", t)
		timeframe, _ = time.ParseDuration(t)
	}

	// Check if limiter store type has been reconfigured by ENV
	storeType := RATE_LIMITER_STORE_MEMORY
	sT := os.Getenv("RATE_LIMITER_STORE_TYPE")

	// Override default config with ENV configuration
	if sT != "" {
		logger.Printf("Read limiter type from RATE_LIMITER_STORE_TYPE env: %s", sT)
		storeType = sT
	}

	config := NewConfig(limit, timeframe, storeType)

	return New(config)
}

func New(config Config) gin.HandlerFunc {
	limiter := newLimiter(config)

	logger.Printf("New Limiter created with rate limit of %d requests per %s", config.Limit, config.Interval)
	return limiter.ApplyRateLimiter()
}

func keyFunc(c *gin.Context) string {
	return c.ClientIP()
}

func errorHandler(c *gin.Context, info ratelimit.Info) {
	c.JSON(http.StatusTooManyRequests, gin.H{"status": http.StatusTooManyRequests, "message": "Too many requests. Try again in " + time.Until(info.ResetTime).String()})
}
