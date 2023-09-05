package controller

import (
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"

	"hackattic-basic-face-detection/app/services"
	"hackattic-basic-face-detection/infra/aws"
	"hackattic-basic-face-detection/infra/hackattic"
)

func SolveProblem(context *gin.Context) {
    session, err := aws.NewSession(
        os.Getenv("AWS_ACCESS_KEY_ID"),
        os.Getenv("AWS_SECRET_ACCESS_KEY"),
        os.Getenv("AWS_REGION"),
    )
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    s3Api := s3.New(session)
    s3Client := aws.NewS3Client(s3Api)
    rekognitionApi := rekognition.New(session)
    rekognitionClient := aws.NewRekognitionClient(rekognitionApi)
    detectFaces := services.NewDetectFaces(s3Client, rekognitionClient)

    hackatticClient := hackattic.HackatticClient{
        AccessToken: os.Getenv("HACKATTIC_ACCESS_TOKEN"),
    }
    problem, err := hackatticClient.GetProblem()

    detectedFaces, err := detectFaces.Perform(problem)
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    tiles, err := services.FindFaceContainingTiles(detectedFaces)
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    result, err := hackatticClient.SubmitSolution(tiles)
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    context.IndentedJSON(http.StatusOK, gin.H{"data": result})
}
