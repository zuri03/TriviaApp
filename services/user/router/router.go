package router

import (
	"context"
	"net/http"
	"os"

	//"github.com/nitishm/go-rejson/v4"
	"github.com/go-redis/redis/v8"
	"github.com/zuri03/user/handlers"
)

func InitRouter(redisJson *redis.Client, signaler chan os.Signal, ctx context.Context) *http.ServeMux {

	router := http.NewServeMux()

	createHandler := handlers.CreateHandler{RedisHandler: redisJson, Signaler: signaler, Ctx: ctx}
	getHandler := handlers.GetHandler{RedisHandler: redisJson, Signaler: signaler, Ctx: ctx}
	/*
		getHandler := handlers.GetHandler{RedisHandler: redisJson, Signaler: signaler}
		deleteHandler := handlers.DeleteHandler{RedisHandler: redisJson, Signaler: signaler}
		updateHandler := handlers.UpdateHandler{RedisHandler: redisJson, Signaler: signaler}
	*/

	router.HandleFunc("/user", func(writer http.ResponseWriter, req *http.Request) {

		writer.Header().Add("Access-Control-Allow-Origin", "*")
		writer.Header().Add("Access-Control-Allow-Methods", "OPTIONS, GET, POST, DELETE, PUT")
		writer.Header().Add("Access-Control-Allow-Headers", "X-PINGOTHER, Content-Type")
		writer.Header().Add("Access-Control-Max-Age", "86400")

		if req.Header.Get("Content-Type") != "" {
			value := req.Header.Values("Content-Type")
			if value[0] != "application/json" {
				writer.WriteHeader(http.StatusBadRequest)
				writer.Write([]byte("Error: expected json"))
				return
			}
		}

		switch req.Method {
		case http.MethodPost:
			createHandler.ServeHTTP(writer, req)
		case http.MethodGet:
			getHandler.ServeHTTP(writer, req)
		/*
			case http.MethodDelete:
				deleteHandler.ServeHTTP(writer, req)
			case http.MethodPut:
				updateHandler.ServeHTTP(writer, req)
			case http.MethodOptions:
				writer.WriteHeader(http.StatusOK)
			case http.MethodGet:
				getHandler.ServeHTTP(writer, req)
		*/
		default:
			writer.WriteHeader(http.StatusMethodNotAllowed)
			writer.Write([]byte("Method not allowed"))
		}
	})

	return router
}
