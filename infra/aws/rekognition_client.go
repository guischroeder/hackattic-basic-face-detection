package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
)

type RekognitionClient struct {
    rekognition *rekognition.Rekognition
}

func NewRekognitionClient(session *session.Session) *RekognitionClient {
    return &RekognitionClient{
        rekognition: rekognition.New(session),
    }
}

type S3Object struct {
    Bucket string
    Name string
}

func (r *RekognitionClient) DetectFaces(s3Object S3Object) (*rekognition.DetectFacesOutput, error) {
    faces, err := r.rekognition.DetectFaces(&rekognition.DetectFacesInput{
        Image: &rekognition.Image{
            S3Object: &rekognition.S3Object{
                Bucket: aws.String(s3Object.Bucket),
                Name: aws.String(s3Object.Name),
            },
        },
    })

    if err != nil {
        return nil, err
    }

    return faces, nil
}
