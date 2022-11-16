package redis

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/alexm24/golang/internal/models"
)

func NewRedisPool(cfg models.RedisConfig) (*redis.Pool, error) {
	redisUrl := fmt.Sprintf("%s:%s", cfg.Url, cfg.Port)
	rp := &redis.Pool{
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisUrl)
			if err != nil {
				return nil, err
			}
			return c, err
		},
	}

	return rp, nil
}
