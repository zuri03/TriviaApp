package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"
	"github.com/zuri03/user/handlers"
)

func InitServer() {

	redisJsonHanlder := rejson.NewReJSONHandler()

	fmt.Println("Connecting to reids")
	redisDBClient := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	redisJsonHanlder.SetGoRedisClient(redisDBClient)
	fmt.Println("Connection successfull")
	http.Handle("/user", &handlers.CreateHandler{RedisHandler: redisJsonHanlder})
	go func() {
		err := http.ListenAndServe(":8081", nil)
		if err != nil {
			log.Fatalf("Error occured initializing server: %s\n", err.Error())
			return
		}
		fmt.Printf("Listening on port 8080 \n")
	}()
}
