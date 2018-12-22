package redis

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/tada3/triton/config"
)

const (
	dataSourceNameFmt = "%s:%d"
)

var client *redis.Client

func init() {
	var err error
	dsn := getDataSourceName()
	fmt.Printf("INFO Connecting to Redis(%s)..", dsn)
	client, err = newRedisClient(dsn, getPassword())
	if err != nil {
		panic(err)
	}
}

func Get(key string) (string, bool) {
	v, err := client.Get(key).Result()
	if err != nil {
		if err != redis.Nil {
			// real error
			fmt.Printf("LOG Redis Get(%s) error: %s", key, err.Error())
		}
		return "", false
	}
	return v, true
}

func Set(k, v string, e time.Duration) {
	err := client.Set(k, v, e).Err()
	if err != nil {
		fmt.Printf("LOG Redis Set(%s, %s, %d) error: %s", k, v, e, err.Error())
	}
}

func newRedisClient(addr string, passwd string) (*redis.Client, error) {
	if addr == "" {
		return nil, errors.New("addr should not be empty string")
	}
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: passwd,
		PoolSize: 30,
	}), nil
}

func getDataSourceName() string {
	cfg := config.GetConfig()
	// Assume cfg is never nil
	return fmt.Sprintf(dataSourceNameFmt,
		cfg.RedisHost, cfg.RedisPort)
}

func getPassword() string {
	cfg := config.GetConfig()
	return cfg.RedisPasswd
}
