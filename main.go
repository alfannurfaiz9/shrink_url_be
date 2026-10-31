package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/alfannurfaiz9/shrink_url_be/internal/config"
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

	dbConfig := config.NewPsqlDB(
		os.Getenv("DBUSER"),
		os.Getenv("DBPASSWORD"),
		os.Getenv("DBHOST"),
		os.Getenv("DBPORT"),
		os.Getenv("DBNAME"),
	)

	pool, dbErr := dbConfig.Connect()
	if dbErr != nil {
		log.Fatalf("Cannot connect to database... \n%s", dbErr)
		return
	}
	defer pool.Close()

	dbCtx, dbCancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := pool.Ping(dbCtx); err != nil {
		log.Fatalf("Database is not ready... \n%s", err)
	}
	defer dbCancel()
	log.Println("Successfully connect to database...")

	rdbConf := config.NewRdb(
		os.Getenv("REDIS_USER"),
		os.Getenv("REDIS_PASS"),
		os.Getenv("REDIS_HOST"),
		os.Getenv("REDIS_PORT"),
	)

	rdCtx, rdCancel := context.WithTimeout(context.Background(), 5*time.Second)
	rdb := rdbConf.Connect()
	defer rdb.Close()

	_, rdErr := rdb.Ping(rdCtx).Result()
	if rdErr != nil {
		log.Fatalf("Redis is not ready... \n%s", rdErr)
	}
	defer rdCancel()
	log.Println("Successfully connect to redis...")

	r := gin.Default()
	router.InitMainRouter(r, pool)

	r.Run(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")))
}
