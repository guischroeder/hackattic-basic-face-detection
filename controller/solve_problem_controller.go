package controller

import (
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"

	"hackattic-basic-face-detection/libs/aws"
	"hackattic-basic-face-detection/libs/basic_face_detection/services"
	"hackattic-basic-face-detection/libs/hackattic"
)

func SolveProblem(context *gin.Context) {
    session, err := aws.NewSession(
        os.Getenv("AWS_ACCESS_KEY_ID"),
        os.Getenv("AWS_SECRET_ACCESS_KEY"),
        os.Getenv("AWS_REGION"),
    )
    s3Api := s3.New(session)
    s3Client := aws.NewS3Client(s3Api)
    rekognitionApi := rekognition.New(session)
    rekognitionClient := aws.NewRekognitionClient(rekognitionApi)
    hackatticClient := hackattic.HackatticClient{
        AccessToken: os.Getenv("HACKATTIC_ACCESS_TOKEN"),
    }

    problem, err := hackatticClient.GetProblem()
    if err != nil {
        context.Error(err)
    }

    detectFaces := services.NewDetectFaces(s3Client, rekognitionClient)
    detectedFaces, err := detectFaces.Perform(problem)
    if err != nil {
        context.Error(err)
    }

    tiles, err := services.FindFaceContainingTiles(detectedFaces)
    if err != nil {
        context.Error(err)
    }

    result, err := hackatticClient.SubmitSolution(tiles)
    if err != nil {
        context.Error(err)
    }

    if len(context.Errors) > 0 {
        context.IndentedJSON(http.StatusBadRequest, gin.H{"errors": context.Errors})
    } else {
        context.IndentedJSON(http.StatusOK, gin.H{"data": result})
    }
}
