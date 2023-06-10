package redis

import "github.com/gomodule/redigo/redis"

func NewRedisPool(redisAddress string) *redis.Pool {
	pool := redis.Pool(redis.Pool{Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", redisAddress)
	}})
	return &pool
}
