package main

import (
	"fmt"

	"github.com/IldarGaleev/todo-backend-service/config"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println(config.AppConfig.ConfigStr)
}
