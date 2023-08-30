package controller

import (
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"hackattic-basic-face-detection/libs/aws"
	"hackattic-basic-face-detection/libs/aws/rekognition"
	"hackattic-basic-face-detection/libs/aws/s3"
	"hackattic-basic-face-detection/libs/hackattic"
)

func SolveProblem(context *gin.Context) {
    hackatticManager := hackattic.HackatticClient{
        AccessToken: os.Getenv("HACKATTIC_ACCESS_TOKEN"),
    }

    problem, err := hackatticManager.GetProblem()
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    resp, err := http.Get(problem.ImageUrl)
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    defer resp.Body.Close()

    session, err := aws.NewSession(
        os.Getenv("AWS_ACCESS_KEY_ID"),
        os.Getenv("AWS_SECRET_ACCESS_KEY"),
        os.Getenv("AWS_REGION"),
    )
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    s3Client := s3.S3(session)
    data, _ := io.ReadAll(resp.Body)
    reader := strings.NewReader(string(data))

    bucket := os.Getenv("S3_BUCKET_NAME")
    path := "media/faces.jpg"
    err = s3.Upload(s3Client, s3.UploadInput{
        Bucket: bucket,
        Path: path,
        Image: io.ReadSeeker(reader), 
    })
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    rekognitionClient := rekognition.Rekognition(session)
    detectedFaces, err := rekognition.DetectFaces(
        rekognitionClient, rekognition.S3Object{
            Bucket: bucket,
            Name: path,
        },
    )

    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    context.IndentedJSON(http.StatusOK, gin.H{"data": detectedFaces})
}
