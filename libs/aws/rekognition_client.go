package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/aws/aws-sdk-go/service/rekognition/rekognitioniface"
)


type RekognitionClient struct {
    RekognitionApi rekognitioniface.RekognitionAPI
}

func NewRekognitionClient(rekognition *rekognition.Rekognition) *RekognitionClient {
    return &RekognitionClient{
        RekognitionApi: rekognition,
    }
}

type S3Object struct {
    Bucket string
    Name string
}

func (r *RekognitionClient) DetectFaces(s3Object S3Object) (*rekognition.DetectFacesOutput, error) {
    faces, err := r.RekognitionApi.DetectFaces(&rekognition.DetectFacesInput{
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
