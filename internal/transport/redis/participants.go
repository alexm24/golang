package redis

import (
	"encoding/json"

	"github.com/gomodule/redigo/redis"

	"github.com/alexm24/golang/internal/models"
)

type ParticipantsRedis struct {
	redisPool *redis.Pool
}

func NewParticipantsRedis(redisPool *redis.Pool) *ParticipantsRedis {
	return &ParticipantsRedis{redisPool}
}

func (p *ParticipantsRedis) CreateParticipant(channel string, user models.PostParticipant) error {
	data, _ := json.Marshal(user)
	redisCon := p.redisPool.Get()
	defer redisCon.Close()

	_, err := redisCon.Do("HSET", channel, *user.Username, string(data))
	if err != nil {
		return err
	}

	_, err = redisCon.Do("EXPIRE", channel, 432000)
	if err != nil {
		return err
	}

	return nil
}
