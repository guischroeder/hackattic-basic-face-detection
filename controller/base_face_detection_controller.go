package controller

import (
    "os"
    "github.com/gin-gonic/gin"
    "net/http"

    "hackattic-base-face-detection/libs/hackattic"
)

func SolveProblem(context *gin.Context) {
    hackatticManager := hackattic.HackatticManager{
        BaseUrl: "https://hackattic.com/challenges/basic_face_detection",
    }
    accessToken := os.Getenv("HACKATTIC_ACCESS_TOKEN")

    problem, err := hackatticManager.GetProblem(accessToken)
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    context.IndentedJSON(http.StatusOK, gin.H{"data": problem})
}
