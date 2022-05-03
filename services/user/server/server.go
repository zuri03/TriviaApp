package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"
	"github.com/zuri03/user/router"
)

func InstantiateRedisDBClientHandler() *rejson.Handler {
	redisJsonHanlder := rejson.NewReJSONHandler()

	fmt.Println("Connecting to reids")
	redisDBClient := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	redisJsonHanlder.SetGoRedisClient(redisDBClient)
	fmt.Println("Connection successfull")
	return redisJsonHanlder
}

func InitServer(signaler chan os.Signal) {

	rejsonRedisHandler := InstantiateRedisDBClientHandler()

	userRouter := router.InitRouter(rejsonRedisHandler, signaler)

	//http.Handle("/user", &handlers.CreateHandler{RedisHandler: rejsonRedisHandler})
	go func() {
		err := http.ListenAndServe(":8081", userRouter)
		if err != nil {
			log.Fatalf("Error occured initializing server: %s\n", err.Error())
			return
		}
		fmt.Printf("Listening on port 8081 \n")
	}()

}
