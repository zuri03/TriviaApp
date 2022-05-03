package main

import (
	"fmt"
	"os"

	"github.com/zuri03/user/server"
)

func main() {
	receiver := make(chan os.Signal, 1)
	server.InitServer(receiver)
	<-receiver
	close(receiver)
	fmt.Println("Signal Received")
	fmt.Println("Shutting down")
}
