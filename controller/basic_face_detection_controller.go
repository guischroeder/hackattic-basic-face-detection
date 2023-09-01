package controller

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"hackattic-basic-face-detection/libs/aws"
	"hackattic-basic-face-detection/libs/aws/rekognitionclient"
	"hackattic-basic-face-detection/libs/aws/s3client"
	"hackattic-basic-face-detection/libs/basicfacedetection"
	"hackattic-basic-face-detection/libs/hackattic"
)

func SolveProblem(context *gin.Context) {
    hackatticClient := hackattic.HackatticClient{
        AccessToken: os.Getenv("HACKATTIC_ACCESS_TOKEN"),
    }
    session, err := aws.NewSession(
        os.Getenv("AWS_ACCESS_KEY_ID"),
        os.Getenv("AWS_SECRET_ACCESS_KEY"),
        os.Getenv("AWS_REGION"),
    )
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    s3Client := s3client.S3(session)
    rekognitionClient := rekognitionclient.Rekognition(session)
    solveProblem := basicfacedetection.NewSolveProblemService(
        hackatticClient,
        s3Client,
        rekognitionClient,
    )

    faceContainingTiles, err := solveProblem.Perform()
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    result, err := hackatticClient.SubmitSolution(faceContainingTiles)
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    context.IndentedJSON(http.StatusOK, gin.H{"data": result})
}
