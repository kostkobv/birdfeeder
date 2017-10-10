package main

import "api"
import "utils"

func main() {
	api.InitServer(":8081", utils.Init()).Start()
}
