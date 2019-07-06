package redis

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/tada3/triton/config"
	"github.com/tada3/triton/logging"
)

const (
	dataSourceNameFmt = "%s:%d"
)

var (
	log    *logging.Entry
	client *redis.Client
)

func init() {
	log = logging.NewEntry("redis")
	var err error
	dsn := getDataSourceName()
	log.Info("Connecting to Redis(%s)..", dsn)
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
			//fmt.Printf("LOG Redis Get(%s) error: %s", key, err.Error())
			log.Error("Redis Get(%s) failed!", key, err)
		}
		return "", false
	}
	return v, true
}

func Set(k, v string, e time.Duration) {
	err := client.Set(k, v, e).Err()
	if err != nil {
		log.Error("Redis Set(%s, %s, %d) failed!", k, v, e, err)
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
