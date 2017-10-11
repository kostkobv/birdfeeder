package main

import (
	"api"
	"utils"
)

func main() {
	api.InitServer(":8081", utils.Init()).Start()
}
