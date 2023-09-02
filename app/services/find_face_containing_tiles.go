package services

import (
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/service/rekognition"

	"hackattic-basic-face-detection/app/helpers"
	"hackattic-basic-face-detection/infra/aws"
	"hackattic-basic-face-detection/infra/hackattic"
	"hackattic-basic-face-detection/infra/httpclient"
)

const IMAGE_HEIGHT = 800.0
const IMAGE_WIDTH = 800.0
const PATH = "media/faces.jpg"

type HackatticClient interface {
    GetProblem() (hackattic.Problem, error)
}
type FindFaceContainingTiles struct {
    hackattic HackatticClient
    s3Client *aws.S3Client
    rekognitionClient *aws.RekognitionClient
}

func NewFindFaceContainingTiles(hackattic HackatticClient,
    s3Client *aws.S3Client, rekognitionClient *aws.RekognitionClient) *FindFaceContainingTiles {
    return &FindFaceContainingTiles{
        hackattic: hackattic,
        s3Client: s3Client,
        rekognitionClient: rekognitionClient,
    }
}

func (f *FindFaceContainingTiles) Perform() ([][2]int, error) {
    problem, err := f.hackattic.GetProblem()
    if err != nil {
        return nil, err
    }

    bucket := os.Getenv("S3_BUCKET_NAME")
    err = f.uploadImageFromUrlToS3(problem.ImageUrl, bucket, PATH)
    if err != nil {
        return nil, err
    }

    detectedFaces, err := f.rekognitionClient.DetectFaces(aws.S3Object{
        Bucket: bucket,
        Name: PATH,
    })
    if err != nil {
        return nil, err
    }

    return f.findFaceContainingTiles(detectedFaces), nil
}

func (f *FindFaceContainingTiles) uploadImageFromUrlToS3(
    imageUrl string, bucket string, path string) error {
    resp, err := httpclient.Client.Get(imageUrl)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    data, _ := io.ReadAll(resp.Body)
    reader := strings.NewReader(string(data))

    err = f.s3Client.Upload(aws.UploadInput{
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

func (f FindFaceContainingTiles) findFaceContainingTiles(
    detectedFaces *rekognition.DetectFacesOutput) [][2]int {
    faceDetails := detectedFaces.FaceDetails
 
    scaledBoundingBoxes := make([]BoundingBox, 0, cap(faceDetails))
    for _, faceDetail := range faceDetails {
        scaledBoundingBox := BoundingBox{
            FaceHeight: *faceDetail.BoundingBox.Height * IMAGE_HEIGHT,
            LeftCoordinate: *faceDetail.BoundingBox.Left * IMAGE_WIDTH,
            TopCoordinate: *faceDetail.BoundingBox.Top * IMAGE_HEIGHT,
            FaceWidth: *faceDetail.BoundingBox.Width * IMAGE_WIDTH,
        }

        scaledBoundingBoxes = append(scaledBoundingBoxes, scaledBoundingBox)
    }

    faceContainingTiles := make([][2]int, 0, cap(scaledBoundingBoxes))
    for _, boundingBox := range scaledBoundingBoxes {
        x := helpers.TilePosition(
            helpers.CenterOfSquare(boundingBox.LeftCoordinate, boundingBox.FaceWidth),
            IMAGE_HEIGHT,
        )
        y := helpers.TilePosition(
            helpers.CenterOfSquare(boundingBox.TopCoordinate, boundingBox.FaceHeight),
            IMAGE_WIDTH,
        )

        faceContainingTiles = append(faceContainingTiles, [2]int{x, y})
    }
    
    return faceContainingTiles
} 
