package main

import (
	"api"
	"config"
	"external"
	"queue"
	"utils"
)

func main() {
	mb := external.InitMessageBirdClient(config.MessageBirdKey)
	q := queue.InitQueue(mb)
	udh := utils.InitEncoder()
	v := utils.InitValidator()
	api.InitServer(config.ServerAddress, v, udh, q).Start()
}
