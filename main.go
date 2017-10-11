package main

import (
	"api"
	"queue"
	"utils"
)

func main() {
	q := queue.InitQueue()
	api.InitServer(":8081", utils.Init(), q).Start()
}
