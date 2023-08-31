package services

import (
	"io"
	"net/http"
	"os"
	"strings"

	"hackattic-basic-face-detection/infra/aws"
	"hackattic-basic-face-detection/app/helpers"
	"hackattic-basic-face-detection/infra/hackattic"

	"github.com/aws/aws-sdk-go/service/rekognition"
)

type HackatticClient interface {
    GetProblem() (hackattic.Problem, error)
}
type SolveProblemService struct {
    hackattic HackatticClient
    s3Client *aws.S3Client
    rekognitionClient *aws.RekognitionClient
}

func NewSolveProblemService(hackattic HackatticClient,
    s3Client *aws.S3Client, rekognitionClient *aws.RekognitionClient) *SolveProblemService {
    return &SolveProblemService{
        hackattic: hackattic,
        s3Client: s3Client,
        rekognitionClient: rekognitionClient,
    }
}

func (s *SolveProblemService) Perform() ([][2]int, error) {
    problem, err := s.hackattic.GetProblem()
    if err != nil {
        return nil, err
    }

    bucket := os.Getenv("S3_BUCKET_NAME")
    path := "media/faces.jpg"

    err = s.uploadImageFromUrlToS3(problem.ImageUrl, bucket, path)
    if err != nil {
        return nil, err
    }

    detectedFaces, err := s.rekognitionClient.DetectFaces(aws.S3Object{
        Bucket: bucket,
        Name: path,
    })
    if err != nil {
        return nil, err
    }

    return s.findFaceContainingTiles(*detectedFaces), nil
}

func (s *SolveProblemService) uploadImageFromUrlToS3(
    imageUrl string, bucket string, path string) error {
    resp, err := http.Get(imageUrl)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    data, _ := io.ReadAll(resp.Body)
    reader := strings.NewReader(string(data))

    err = s.s3Client.Upload(aws.UploadInput{
        Bucket: bucket,
        Path: path,
        Image: io.ReadSeeker(reader), 
    })
    if err != nil {
        return err
    }

    return nil
}

type BoundingBox struct {
    FaceHeight float64
    LeftCoordinate float64
    TopCoordinate float64
    FaceWidth float64
}

func (s SolveProblemService) findFaceContainingTiles(
    detectedFaces rekognition.DetectFacesOutput) [][2]int {
    faceDetails := detectedFaces.FaceDetails

    imageHeight := 800.0
    imageWidth := 800.0
    scaledBoundingBoxes := make([]BoundingBox, 0, cap(faceDetails))
    for _, faceDetail := range faceDetails {
        scaledBoundingBox := BoundingBox{
            FaceHeight: *faceDetail.BoundingBox.Height * imageHeight,
            LeftCoordinate: *faceDetail.BoundingBox.Left * imageWidth,
            TopCoordinate: *faceDetail.BoundingBox.Top * imageHeight,
            FaceWidth: *faceDetail.BoundingBox.Width * imageWidth,
        }

        scaledBoundingBoxes = append(scaledBoundingBoxes, scaledBoundingBox)
    }

    faceContainingTiles := make([][2]int, 0, cap(scaledBoundingBoxes))
    for _, boundingBox := range scaledBoundingBoxes {
        x := helpers.TilePosition(
            helpers.CenterOfSquare(boundingBox.LeftCoordinate, boundingBox.FaceWidth),
            imageHeight,
        )
        y := helpers.TilePosition(
            helpers.CenterOfSquare(boundingBox.TopCoordinate, boundingBox.FaceHeight),
            imageWidth,
        )

        faceContainingTiles = append(faceContainingTiles, [2]int{x, y})
    }
    
    return faceContainingTiles
} 
