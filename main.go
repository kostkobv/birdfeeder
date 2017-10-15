package main

import (
	"api"
	"queue"
	"utils"
)

func main() {
	q := queue.InitQueue()
	udh := utils.InitEncoder()
	v := utils.InitValidator()
	api.InitServer(":8081", v, udh, q).Start()
}
