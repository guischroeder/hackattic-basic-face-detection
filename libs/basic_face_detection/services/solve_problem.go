package services

import (
	"os"

	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/aws/aws-sdk-go/service/s3"

	"hackattic-basic-face-detection/libs/aws"
	"hackattic-basic-face-detection/libs/hackattic"
)

func SolveProblem() (hackattic.Result, error) {
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
        return hackattic.Result{}, err
    }

    detectFaces := NewDetectFaces(s3Client, rekognitionClient)
    detectedFaces, err := detectFaces.Perform(problem)
    if err != nil {
        return hackattic.Result{}, err
    }

    tiles, err := FindFaceContainingTiles(detectedFaces)
    if err != nil {
        return hackattic.Result{}, err
    }

    result, err := hackatticClient.SubmitSolution(tiles)
    if err != nil {
        return hackattic.Result{}, err
    }

    return result, nil
}
