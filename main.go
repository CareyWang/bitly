package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"io/ioutil"
	"net/http"
	"time"
)

type Response struct {
	Code     int
	Message  string
	LongUrl  string
	ShortUrl string
}

type Bitly struct {
	CreatedAt      string        `json:"created_at"`
	ID             string        `json:"id"`
	Link           string        `json:"link"`
	CustomBitlinks []interface{} `json:"custom_bitlinks"`
	LongURL        string        `json:"long_url"`
	Archived       bool          `json:"archived"`
	Tags           []interface{} `json:"tags"`
	Deeplinks      []interface{} `json:"deeplinks"`
	References     struct {
		Group string `json:"group"`
	} `json:"references"`
}

// api文档：https://dev.bitly.com/v4_documentation.html
const endPoint string = "https://api-ssl.bitly.com/v4/shorten"
const defaultToken string = ""
const defaultPort int = 8001

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	token := flag.String("token", "", "Bitly api token")
	port := flag.Int("port", defaultPort, "服务端口")
	cache := flag.Int("cache", 0, "是否使用 redis 缓存")

	flag.Parse()

	if *token == "" {
		panic("Bitly api token is required.")
	}

	router.GET("/", func(c *gin.Context) {
		res := &Response{
			Code:     0,
			Message:  "",
			LongUrl:  "",
			ShortUrl: "",
		}

		longUrl := c.Query("longUrl")
		if longUrl == "" {
			res.Message = "longUrl为空"
			c.JSON(400, *res)
			return
		}

		_longUrl, _ := base64.StdEncoding.DecodeString(longUrl)
		longUrl = string(_longUrl)
		res.LongUrl = longUrl

		// 先查缓存
		redisClient := initRedis()
		if 1 == *cache {
			_shortUrl, err := redisClient.Get(longUrl).Result()
			if err == nil && _shortUrl != "" {
				res.Code = 1
				res.ShortUrl = _shortUrl
				c.JSON(200, res)
				return
			}
		}

		contentType := "application/json"
		postField := fmt.Sprintf("{\"long_url\":\"%s\"}", longUrl)
		req, err := http.NewRequest("POST", endPoint, bytes.NewBuffer([]byte(postField)))
		if err != nil {
			res.Message = err.Error()
			c.JSON(500, *res)
			return
		}

		req.Header.Add("Content-Type", contentType)
		req.Header.Add("Accept", contentType)
		req.Header.Add("Host", "api-ssl.bitly.com")
		req.Header.Add("Authorization", "Bearer "+*token)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			res.Message = err.Error()
			c.JSON(500, *res)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		var bitly Bitly
		err = json.Unmarshal(body, &bitly)

		// 写缓存
		res.ShortUrl = bitly.Link
		if 1 == *cache {
			ttl, err := time.ParseDuration("300s")
			if err != nil {
				res.Message = err.Error()
				c.JSON(500, *res)
				return
			}

			redisClient.Set(res.LongUrl, res.ShortUrl, ttl)
		}

		res.Code = 1
		c.JSON(200, res)
	})

	router.Run(fmt.Sprintf(":%d", *port))
}

func initRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
