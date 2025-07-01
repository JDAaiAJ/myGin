package main

import (
	"employee-salary-gin/config"
	"employee-salary-gin/routers"
	"employee-salary-gin/utils"
)

func main() {
	config.InitDB()
	utils.InitRedis()

	r := routers.SetupRouter()
	r.Run(":8080")
}
