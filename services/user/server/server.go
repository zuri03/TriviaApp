package server

import (
	"fmt"
	"log"
	"net/http"
)

func InitServer() {

	http.HandleFunc("/User", func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.Header().Set("Access-Control-Allow-Origin", "*")

	})

	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatalf("Error occured initializing server: %s\n", err.Error())
			return
		}
		fmt.Printf("Listening on port 8080 \n")
	}()
}
