package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

var pool *redis.Pool

func initPool(address string, maxIdel, MaxActive int, IdleTimeout time.Duration) {
	pool = &redis.Pool{
		MaxActive:   MaxActive,
		MaxIdle:     maxIdel,
		IdleTimeout: IdleTimeout,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", address)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", "tengt"); err != nil {
				c.Close()
				return nil, err
			}
			fmt.Println("连接成功....")
			return c, nil
		},
	}
}
