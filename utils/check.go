package utils

import "os"

func CheckEnvironment() {
	if os.Getenv("APP_NAME") == "" {
		panic("APP_NAME is not set")
	}
	if os.Getenv("BOT_TOKEN") == "" {
		panic("BOT_TOKEN is not set")
	}
	if os.Getenv("REDIS_HOST") == "" {
		panic("REDIS_HOST is not set")
	}
	if os.Getenv("REDIS_PORT") == "" {
		panic("REDIS_PORT is not set")
	}
	if os.Getenv("REDIS_PASSWORD") == "" {
		panic("REDIS_PASSWORD is not set")
	}
	if os.Getenv("REDIS_USER") == "" {
		panic("REDIS_USER is not set")
	}
	if os.Getenv("BOT_MSG_API") == "" {
		panic("BOT_MSG_API is not set")
	}
	if os.Getenv("TWITTER_USER") == "" {
		panic("TWITTER_USER is not set")
	}
	if os.Getenv("TWITTER_PASSWORD") == "" {
		panic("TWITTER_PASSWORD is not set")
	}
	if os.Getenv("TARGET_USER_ID") == "" {
		panic("TWITTER_USER_ID is not set")
	}
}
