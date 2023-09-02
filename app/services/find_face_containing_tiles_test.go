package services

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"hackattic-basic-face-detection/infra/aws"
	"hackattic-basic-face-detection/infra/hackattic"
	"hackattic-basic-face-detection/infra/httpclient"

	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/aws/aws-sdk-go/service/rekognition/rekognitioniface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/stretchr/testify/assert"
)

var (
	getRequestFunc func(url string) (*http.Response, error)
)
type httpClientGetMock struct {}
func (h *httpClientGetMock) Get(request string) (*http.Response, error) {
    return getRequestFunc(request)
}

type hackatticClientMock struct {}
func (h hackatticClientMock) GetProblem() (hackattic.Problem, error) {
    return hackattic.Problem{
        ImageUrl: "http://image.com",
    }, nil
}

type s3ApiMock struct {
    called bool
    s3iface.S3API
}
func (s *s3ApiMock) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
    s.called = true
    return &s3.PutObjectOutput{}, nil
}

type rekognitionApiMock struct {
    called bool
    rekognitioniface.RekognitionAPI
}
func (r *rekognitionApiMock) DetectFaces(s3Object *rekognition.DetectFacesInput) (*rekognition.DetectFacesOutput, error) {
    r.called = true

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

    return &rekognition.DetectFacesOutput{
        FaceDetails: []*rekognition.FaceDetail{faceDetail, faceDetail, faceDetail},
    }, nil
}

func TestFindFaceContainingTiles(t *testing.T) {
    getRequestFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("success")),
		}, nil
	}
    httpclient.Client = &httpClientGetMock{}
    hackattic := hackatticClientMock{}

    s3ApiMock := &s3ApiMock{}
    s3Client := &aws.S3Client{
        S3Api: s3ApiMock,
    }
    rekognitionApiMock := &rekognitionApiMock{}
    rekognitionClient := &aws.RekognitionClient{
        RekognitionApi: rekognitionApiMock,
    }

    findFaceContainingTiles := NewFindFaceContainingTiles(
        hackattic,
        s3Client,
        rekognitionClient,
    )

    result, _ := findFaceContainingTiles.Perform()

    if s3ApiMock.called != true {
        t.Errorf("Expected PutObject method to be called")
    }
    if rekognitionApiMock.called != true {
        t.Errorf("Expected DetectFaces method to be called")
    }

    assert.Equal(t, result, [][2]int{{2,2}, {2,2},{2,2}})
}
