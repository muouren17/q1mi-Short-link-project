package sequence

import (
	"fmt"

	"github.com/go-redis/redis"
)

type Redis struct {
	client *redis.Client
}

const (
	IncryKey = "IncryCount"
)

func NewRedis(addr string, password string) Sequence {
	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})
	return &Redis{
		client: c,
	}
}

func (r *Redis) Next() (uint64, error) {
	number, err := r.client.Incr(IncryKey).Result()
	if err != nil {
		return 0, err
	}
	if number < 0 {
		return 0, fmt.Errorf("redis sequence %q returned a negative value: %d", IncryKey, number)
	}
	return uint64(number), nil
}
