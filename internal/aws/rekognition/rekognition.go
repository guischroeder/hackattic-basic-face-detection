package rekognition

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/rekognition"
    "github.com/aws/aws-sdk-go/service/rekognition/rekognitioniface"
)

type Rekognition struct {
    Rekognition rekognitioniface.RekognitionAPI
}

func NewRekognition(session *session.Session) *Rekognition {
    return &Rekognition{
        Rekognition: rekognition.New(session),
    }
}

type DetectFacesOutput struct {
    *rekognition.DetectFacesOutput
}

func (r *Rekognition) DetectFaces(bucket string, name string) (*DetectFacesOutput, error) {
    faces, err := r.Rekognition.DetectFaces(&rekognition.DetectFacesInput{
        Image: &rekognition.Image{
            S3Object: &rekognition.S3Object{
                Bucket: aws.String(bucket),
                Name: aws.String(name),
            },
        },
    })

    if err != nil {
        return nil, err
    }

    return &DetectFacesOutput{faces}, nil
}
