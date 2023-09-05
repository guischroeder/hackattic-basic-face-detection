package services

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/stretchr/testify/assert"
)

func TestDetectFaces(t *testing.T) {
    var heigh float64 = 0.5
    var top float64 = 0.05
    var left float64 = 0.05
    var width float64 = 0.5

    faceDetail := &rekognition.FaceDetail{
        BoundingBox: &rekognition.BoundingBox {
            Height: &heigh,
            Top: &top,
            Left: &left,
            Width: &width,
        },
    }

    detectedFaces := &rekognition.DetectFacesOutput{
        FaceDetails: []*rekognition.FaceDetail{faceDetail, faceDetail, faceDetail},
    }
    result, _ := FindFaceContainingTiles(detectedFaces)

    assert.Equal(t, result, [][2]int{{2,2}, {2,2},{2,2}})
}
