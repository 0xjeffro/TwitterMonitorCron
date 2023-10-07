package database

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"os"
)

var conn redis.Conn

func connect() error {
	HOST := os.Getenv("REDIS_HOST")
	PORT := os.Getenv("REDIS_PORT")
	USER := os.Getenv("REDIS_USER")
	PASSWORD := os.Getenv("REDIS_PASSWORD")
	var err error
	conn, err = redis.Dial("tcp", HOST+":"+PORT,
		redis.DialUsername(USER),
		redis.DialPassword(PASSWORD),
		redis.DialUseTLS(true),
	)
	return err
}

func GetConn() redis.Conn {
	if conn != nil {
		return conn
	} else {
		err := connect()
		if err != nil {
			fmt.Println(">>> Get connection failed")
			panic(err)
		} else {
			return conn
		}
	}
}
