package store

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/go-redis/redis"
)

// Redis struct
type Redis struct {
	client *redis.Client
}

// NewRedis function
func NewRedis() *Redis {

	// Addr: ":6379",
	// Addr: os.Getenv("REDIS_URL"),
	c := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URL"),
	})
	pong, err := c.Ping().Result()
	if err != nil {
		panic(err)
	}
	log.Println("Pong", pong)
	return &Redis{client: c}
}

// Save an object to cache
func (r *Redis) Save(urlVal string, value string) error {
	// c := redis.NewClient(&redis.Options{
	// 	Addr: os.Getenv("REDIS_URL"),
	// })
	uri, err := url.Parse(urlVal)
	if err != nil {
		return err
	}
	key := uri.Hostname()
	field := strings.Trim(uri.Path, " ")
	if field == "" {
		field = "/"
	}
	if value == "" {
		value = " "
	}
	hm := map[string]interface{}{}
	hm[field] = value
	// log.Println("REDIS SAVE: ", key, field, value)
	return r.client.HMSet(key, hm).Err()
}

// Fetch an object from cache
func (r *Redis) Fetch(key string) (string, error) {
	// c := redis.NewClient(&redis.Options{
	// 	Addr: os.Getenv("REDIS_URL"),
	// })
	uri, err := url.Parse(key)
	if err != nil {
		return "", err
	}
	if uri.Hostname() == "" {
		log.Println(key)
		return "", fmt.Errorf("unknown hostname")
	}
	field := strings.Trim(uri.Path, " ")
	if field == "" {
		field = "/"
	}
	res := r.client.HGet(uri.Hostname(), field)
	return res.Val(), nil
}
