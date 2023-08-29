package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"

    "hackattic-base-face-detection/controller"
)

func main() {
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatalf("Error loading .env file: %s", err)
    }

    router := gin.Default()
    router.POST("/solve", controller.SolveProblem)
    router.Run()
}

