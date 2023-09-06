package main

import (
	"log"

	"github.com/joho/godotenv"

    "hackattic-basic-face-detection/internal/basicfacedetection"
)

func main() {
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatalf("Error loading .env file: %s", err)
    }

    result, err := basicfacedetection.SolveProblem()
    if err != nil {
        log.Fatalln(err)
    }

    log.Println(result)
}

