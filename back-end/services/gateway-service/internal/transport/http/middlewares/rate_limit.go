package middlewares

import (
	"context"
	"errors"
	"fmt"
	applicationServices "gateway/internal/domain/application-services"
	httpErrors "gateway/pkg/errors/http"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/sethvargo/go-limiter"
	"github.com/sethvargo/go-limiter/httplimit"
	"github.com/sethvargo/go-redisstore"
)

type RateLimiter interface {
	Apply(tokensLimit int) func(c *gin.Context)
}

type rateLimiter struct {
	redisPool *redis.Pool
}

func NewRateLimiter(pool *redis.Pool) *rateLimiter {
	return &rateLimiter{
		redisPool: pool,
	}
}

func (r rateLimiter) Apply(tokensLimit int) func(*gin.Context) {
	// create in-memory store
	store, err := redisstore.NewWithPool(&redisstore.Config{
		Tokens:   uint64(tokensLimit),
		Interval: time.Minute,
	}, r.redisPool)

	if err != nil {
		panic(err)
	}

	return func(c *gin.Context) {
		userFromContext, exists := c.Get("user")
		var userID string
		if exists {
			user := userFromContext.(applicationServices.User)
			userID = user.ID
		}
		err := limit(store, userID, c)

		if err != nil {
			httpErrors.BadRequest(c, err.Error())
			return
		}
		c.Next()
	}
}

func limit(store limiter.Store, userID string, c *gin.Context) error {
	path := c.Request.URL.EscapedPath()
	clientIP := c.ClientIP()
	key := fmt.Sprintf("%s/%s/%s", path, userID, clientIP)
	// Take from the store.
	limit, remaining, reset, ok, err := store.Take(context.Background(), key)
	if err != nil {
		return err
	}

	resetTime := int(math.Floor(time.Unix(0, int64(reset)).UTC().Sub(time.Now().UTC()).Seconds()))
	log.Println(resetTime)

	// Set headers (we do this regardless of whether the request is permitted).
	c.Writer.Header().Set(httplimit.HeaderRateLimitLimit, strconv.FormatUint(limit, 10))
	c.Writer.Header().Set(httplimit.HeaderRateLimitRemaining, strconv.FormatUint(remaining, 10))
	c.Writer.Header().Set(httplimit.HeaderRateLimitReset, strconv.Itoa(resetTime))

	// Fail if there were no tokens remaining.
	if !ok {
		c.Writer.Header().Set(httplimit.HeaderRetryAfter, strconv.Itoa(resetTime))
		return errors.New("too many attempts")
	}
	return nil
}
