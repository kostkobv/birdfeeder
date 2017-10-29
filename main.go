package main

import (
	"api"
	"config"
	"external"
	"fmt"
	"queue"
	"utils"
)

func main() {
	mb := external.InitMessageBirdClient(config.MessageBirdKey)
	q := queue.InitQueue(mb)
	udh := utils.InitEncoder()
	v := utils.InitValidator()
	err := api.InitServer(config.ServerAddress, v, udh, q).Start()
	fmt.Println(err)
}
