package main

import (
	"FileStore-Server/config"
	"FileStore-Server/service/Microservice/apigw/route"
)

func main() {
	r := route.Router()
	r.Run(config.UploadServiceHost)
}