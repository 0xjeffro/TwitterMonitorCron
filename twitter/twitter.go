package twitter

import (
	"TwitterMonitorCron/database"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	twitterscraper "github.com/n0madic/twitter-scraper"
	"io"
	"log"
	"net/http"
	"os"
)

func login(scraper *twitterscraper.Scraper) error {
	c := database.GetConn()
	appName := os.Getenv("APP_NAME")
	// 查找redis中是否存在cookies
	cacheCookies, err := c.Do("GET", appName+"_"+os.Getenv("TWITTER_USER")+"_cookies")
	if err != nil {
		panic(err)
	}
	if cacheCookies != nil {
		// convert cookies string to []*http.Cookie
		var cookies []*http.Cookie
		err := json.NewDecoder(bytes.NewReader(cacheCookies.([]byte))).Decode(&cookies)
		if err != nil {
			fmt.Println(">>> Cookies found but can't convert to []*http.Cookie")
			return err
		}
		// load cookies
		scraper.SetCookies(cookies)

		if scraper.IsLoggedIn() {
			fmt.Println(">>> Cookies found and logged in successfully")
			return nil
		} else {
			err := scraper.Login(os.Getenv("TWITTER_USER"), os.Getenv("TWITTER_PASSWORD"))
			if err != nil {
				return err
			} else {
				fmt.Println(">>> Cookies found but may be expired, password logged in successfully")
				fmt.Println(">>> Saving cookies to redis")
				cookies := scraper.GetCookies()
				// serialize to JSON
				js, _ := json.Marshal(cookies)
				// convert to string
				cookiesString := string(js)
				// save to redis
				_, err := c.Do("SET", appName+"_"+os.Getenv("TWITTER_USER")+"_cookies", cookiesString)
				if err != nil {
					fmt.Println(">>> Saving cookies to redis failed")
					return err
				}
				return nil
			}
		}
	} else {
		err := scraper.Login(os.Getenv("TWITTER_USER"), os.Getenv("TWITTER_PASSWORD"))
		if err != nil {
			fmt.Println(">>> No cookies found, password logged in failed")
			return err
		} else {
			fmt.Println(">>> No cookies found, password logged in successfully")
			fmt.Println(">>> Saving cookies to redis")
			cookies := scraper.GetCookies()
			// serialize to JSON
			js, _ := json.Marshal(cookies)
			// convert to string
			cookiesString := string(js)
			// save to redis
			_, err := c.Do("SET", appName+"_"+os.Getenv("TWITTER_USER")+"_cookies", cookiesString)
			if err != nil {
				fmt.Println(">>> Saving cookies to redis failed")
				return err
			}
			return nil
		}
	}
}

func GetTweets() {
	scraper := twitterscraper.New().WithReplies(true)

	err := login(scraper)
	if err != nil {
		fmt.Println(">>> Twitter login failed")
		panic(err)
	}
	scraper.SetSearchMode(twitterscraper.SearchLatest)
	for tweet := range scraper.SearchTweets(context.Background(), "from:"+os.Getenv("TARGET_USER_ID"), 5) {
		if tweet.Error != nil {
			panic(tweet.Error)
		}
		Pin := false
		action := "🔫 "
		if tweet.IsRetweet {
			action = "🔫 "
		} else if tweet.IsReply {
			action = "💬 "
		} else if tweet.IsQuoted {
			action = "🔫 "
			Pin = true
		} else {
			Pin = true
		}

		conn := database.GetConn()
		// 判断是否已经存在
		isMember, err := redis.Bool(
			conn.Do("SISMEMBER",
				os.Getenv("APP_NAME")+"_"+os.Getenv("TARGET_USER_ID")+"_tweets",
				tweet.PermanentURL))
		if err != nil {
			panic(err)
		}
		if !isMember {
			// 不存在则添加
			_, err = conn.Do("SADD",
				os.Getenv("APP_NAME")+"_"+os.Getenv("TARGET_USER_ID")+"_tweets",
				tweet.PermanentURL)
			if err == nil {
				// 发送消息
				Text := action + string([]rune(tweet.Text)[:18])
				if len([]rune(tweet.Text)) > 18 {
					Text += "..."
				}
				URL := tweet.PermanentURL
				Token := os.Getenv("BOT_TOKEN")

				type PostData struct {
					Text       string `json:"text"`
					TwitterURL string `json:"twitter_url"`
					Token      string `json:"token"`
					Pin        bool   `json:"pin"` // 是否置顶
				}
				// 发送post请求
				data := PostData{
					Text:       Text,
					TwitterURL: URL,
					Token:      Token,
					Pin:        Pin,
				}
				bytesData, err := json.Marshal(data)
				if err != nil {
					panic(err)
				}
				reader := bytes.NewReader(bytesData)
				url := os.Getenv("BOT_MSG_API")
				req, err := http.Post(url, "application/json", reader)
				if err != nil {
					log.Panicln(err)
					return
				}
				fmt.Println(action, tweet.ID, tweet.PermanentURL, tweet.Text, tweet.Timestamp, Pin)
				defer func(Body io.ReadCloser) {
					err := Body.Close()
					if err != nil {
						log.Println(err)
					}
				}(req.Body)
			}
		} else {
			// 存在则跳过
			fmt.Println(">>> Tweet already exists " + tweet.PermanentURL)
			continue
		}
	}
}
