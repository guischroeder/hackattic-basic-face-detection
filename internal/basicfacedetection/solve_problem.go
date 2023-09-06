package basicfacedetection

import (
    "errors"
    "math"
    "os"

    "hackattic-basic-face-detection/internal/aws"
    "hackattic-basic-face-detection/internal/hackattic"
    "hackattic-basic-face-detection/internal/aws/rekognition"
    "hackattic-basic-face-detection/internal/aws/s3"
)

func SolveProblem() (hackattic.Result, error) {
    session, err := aws.NewSession(
        os.Getenv("AWS_ACCESS_KEY_ID"),
        os.Getenv("AWS_SECRET_ACCESS_KEY"),
        os.Getenv("AWS_REGION"),
    )

    if err != nil { 
        return hackattic.Result{}, err 
    }

    s3 := s3.NewS3(session)
    rekognitionClient := rekognition.NewRekognition(session)
    hackatticClient := hackattic.HackatticClient{
        AccessToken: os.Getenv("HACKATTIC_ACCESS_TOKEN"),
    }

    problem, err := hackatticClient.GetProblem()
    if err != nil {
        return hackattic.Result{}, err
    }

    detectFaces := NewDetectFaces(s3, rekognitionClient)
    detectedFaces, err := detectFaces.Perform(problem)
    if err != nil {
        return hackattic.Result{}, err
    }

    tiles, err := findFaceContainingTiles(detectedFaces)
    if err != nil {
        return hackattic.Result{}, err
    }

    result, err := hackatticClient.SubmitSolution(tiles)
    if err != nil {
        return hackattic.Result{}, err
    }

    return result, nil
}

const IMAGE_HEIGHT = 800.0
const IMAGE_WIDTH = 800.0

type BoundingBox struct {
    FaceHeight float64
    LeftCoordinate float64
    TopCoordinate float64
    FaceWidth float64
}

func findFaceContainingTiles(detectedFaces *rekognition.DetectFacesOutput) ([][2]int, error) {
    faceDetails := detectedFaces.FaceDetails

    if len(faceDetails) == 0 {
        return nil, errors.New("Error when trying to find tiles containing faces because no faces were detected.")
    }

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
        x := tilePosition(
            centerOfSquare(boundingBox.LeftCoordinate, boundingBox.FaceWidth),
            IMAGE_HEIGHT,
            )
        y := tilePosition(
            centerOfSquare(boundingBox.TopCoordinate, boundingBox.FaceHeight),
            IMAGE_WIDTH,
            )

        faceContainingTiles = append(faceContainingTiles, [2]int{x, y})
    }

    return faceContainingTiles, nil
}

func centerOfSquare(coordinate float64, dimension float64) float64 {
    return coordinate + (dimension / 2)
}

func tilePosition(value float64, imgDimensionSize float64) int {
    numberOfTiles := 8
    tileSize := imgDimensionSize / float64(numberOfTiles)

    position := math.Floor(value / tileSize)

    if v := int(value) % 100; v == 0 {
        position = position - 1
    }

    return int(position)
}
