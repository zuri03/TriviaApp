package router

import (
	"fmt"
	"net/http"

	"github.com/nitishm/go-rejson/v4"
	"github.com/zuri03/user/handlers"
)

func InitRouter(redisJson *rejson.Handler) *http.ServeMux {

	router := http.NewServeMux()

	createHandler := handlers.CreateHandler{RedisHandler: redisJson}
	getHandler := handlers.GetHandler{RedisHandler: redisJson}
	deleteHandler := handlers.DeleteHandler{RedisHandler: redisJson}

	router.HandleFunc("/user", func(writer http.ResponseWriter, req *http.Request) {

		if req.Header.Get("Content-Type") != "" {
			value := req.Header.Values("Content-Type")
			if value[0] != "application/json" {
				writer.WriteHeader(http.StatusBadRequest)
				writer.Write([]byte("Error: expected json"))
				return
			}
		}

		fmt.Println("inside of user router")
		switch req.Method {
		case http.MethodGet:
			getHandler.ServeHTTP(writer, req)
		case http.MethodPost:
			createHandler.ServeHTTP(writer, req)
		case http.MethodDelete:
			deleteHandler.ServeHTTP(writer, req)
		default:
			writer.WriteHeader(http.StatusMethodNotAllowed)
			writer.Write([]byte("Method not allowed"))
		}
	})

	return router
}
