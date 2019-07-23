package main

import (
	"FileStore-Server/config"
	"FileStore-Server/route"
)

func main() {
	// gin框架
	router := route.Router()
	router.Run(config.UploadServiceHost)
}