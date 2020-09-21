package main

import (
	"errors"
	"fmt"
	"redis-golang-cnwnc/cache"
	"redis-golang-cnwnc/connectdb"
	"redis-golang-cnwnc/entities"
	"redis-golang-cnwnc/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
)

var (
	rdb          *redis.Client
	productCache cache.ProductCache
)

func IncRequestCount(key string, rateLimit int64, second int64) error {
	err := rdb.Watch(func(tx *redis.Tx) error {
		_ = tx.SetNX(key, 0, time.Duration(second)*time.Second)
		count, err := tx.Incr(key).Result()
		fmt.Print("số lượt truy cập", count)
		if count > rateLimit {
			err = errors.New("rate limited")
		}
		if err != nil {
			return err
		}
		return nil
	}, key)
	return err
}

func UseRateLimit(rateLimit int64, second int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIp := c.ClientIP()
		key := "RATE_LIMIT_COUNT_" + clientIp
		err := IncRequestCount(key, rateLimit, second)
		if err != nil {
			c.AbortWithStatus(403)
			return
		}
		c.Next()
	}
}

func GetAll(c *gin.Context) {
	db, err := connectdb.GetDB()
	if err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		productModel := models.ProductModel{
			Db: db,
		}
		err1 := UseRateLimit(5, 120)
		if err1 != nil {
			var products *[]entities.Product = productCache.Get("products")
			fmt.Println(products)
		} else {
			procs, err2 := productModel.GetAll()
			productCache.Set("products", &procs)
			if err2 != nil {
				c.AbortWithStatus(404)
				fmt.Println(err)
			} else {
				c.JSON(200, procs)
			}
		}

	}
}

func main() {
	app := gin.Default()
	//test
	// app.GET("/", func(c *gin.Context) {
	// 	c.JSON(200, "Hello World")
	// })
	app.GET("/getall", GetAll)
	rdb = redis.NewClient(&redis.Options{
		Addr:        "localhost:6379", // use default Addr
		Password:    "",               // no password set
		DB:          0,                // use default database
		PoolTimeout: time.Minute,      // since we user transaction so it can take a long time
	})
	app.Run(":8080")
}
