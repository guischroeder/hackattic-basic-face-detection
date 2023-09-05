package services

import (
	"errors"

	"github.com/aws/aws-sdk-go/service/rekognition"

	"hackattic-basic-face-detection/libs/basic_face_detection/helpers"
)

const IMAGE_HEIGHT = 800.0
const IMAGE_WIDTH = 800.0

type BoundingBox struct {
    FaceHeight float64
    LeftCoordinate float64
    TopCoordinate float64
    FaceWidth float64
}

func FindFaceContainingTiles(detectedFaces *rekognition.DetectFacesOutput) ([][2]int, error) {
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
    
    return faceContainingTiles, nil
}
