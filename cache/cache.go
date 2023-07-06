package cache

import (
	"context"
	"os"
	"strconv"
	"time"

	log "github.com/nxs23/utility-library-go/logging/logger"
	"github.com/redis/go-redis/v9"
)

const (
	DEFAULT_CACHE_DB = 0
	DEFULT_CAHCE_TTL = 3600
)

var logger *log.Logger

func init() {
	// Create new Logger
	logger = log.NewLogger("cache")
}

type RedisAdapter struct {
	host     string
	port     string
	db       int
	password string
	ttl      time.Duration
	prefix   string
	client   *redis.Client
}

func Default() *RedisAdapter {
	// Check if host has been configured by ENV
	host := os.Getenv("CACHE_HOST")
	// default host value
	if host == "" {
		host = "localhost"
	}

	// Check if port has been configured by ENV
	port := os.Getenv("CACHE_PORT")
	// default port value
	if port == "" {
		port = "6379"
	}

	// Check if db has been configured by ENV
	db := DEFAULT_CACHE_DB
	dbn := os.Getenv("CACHE_DB")
	// override default db
	if dbn != "" {
		db, _ = strconv.Atoi(dbn)
	}

	// Check if password has been configured by ENV
	password := os.Getenv("CACHE_PASSWORD")
	if password == "" {
		password = ""
	}

	// Check if prefix has been configured by ENV
	prefix := os.Getenv("CACHE_PREFIX")
	if prefix == "" {
		prefix = ""
	}

	// Check if TTL has been configured by ENV
	ttl := DEFULT_CAHCE_TTL
	cttl := os.Getenv("CACHE_TTL")
	if cttl != "" {
		ttl, _ = strconv.Atoi(cttl)
	}

	return NewRedisAdapter(host, port, db, password, time.Duration(ttl), prefix)
}

func NewRedisAdapter(host string, port string, db int, password string, ttl time.Duration, prefix string) *RedisAdapter {
	redisAdapter := &RedisAdapter{
		host:     host,
		port:     port,
		db:       db,
		password: password,
		ttl:      ttl,
		prefix:   prefix,
	}

	client := redisAdapter.NewClient()
	redisAdapter.client = client

	logger.Printf("Created new cache adapter at: %s:%s", host, port)
	return redisAdapter
}

func (adapter *RedisAdapter) NewClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     adapter.host + ":" + adapter.port,
		Password: adapter.password,
		DB:       adapter.db,
	})
}

func (adapter *RedisAdapter) GetClient() *redis.Client {
	return adapter.client
}

func (adapter *RedisAdapter) Ping() error {
	ctx := context.Background()
	_, err := adapter.client.Ping(ctx).Result()
	if err != nil {
		logger.Printf("Error pinging cache at: %s:%s, %s", adapter.host, adapter.port, err.Error())
		return err
	}

	logger.Printf("Successfully pinged cache at: %s:%s", adapter.host, adapter.port)

	return nil
}
