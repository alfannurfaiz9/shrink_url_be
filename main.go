package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alfannurfaiz9/shrink_url_be/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title           			ShrinkURL Backend
// @version         			1.0
// @description     			API docs for ShrinkURL
// @host      					localhost:9000
// @BasePath  					/
// @securityDefinitions.apikey	BearerToken
// @in							header
// @name						Authorization
// @description					Bearer Token used as identity for accessing backend
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err.Error())
		return
	}

	r := gin.Default()
	router.InitMainRouter(r)

	r.Run(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")))
}
