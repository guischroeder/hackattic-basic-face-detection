package rekognitionclient

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/rekognition"
)

type S3Object struct {
    Bucket string
    Name string
}

func DetectFaces(client *rekognition.Rekognition, s3Object S3Object) (*rekognition.DetectFacesOutput, error) {
    faces, err := client.DetectFaces(&rekognition.DetectFacesInput{
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
