package s3

import (
    "io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type UploadInput struct {
    Bucket string
    Path string
    Image io.ReadSeeker
}

func Upload(client *s3.S3, input UploadInput) error {
    putObjectInput := s3.PutObjectInput{
        Bucket: aws.String(input.Bucket),
        Key: aws.String(input.Path),
        Body: input.Image,
    }

    _, err := client.PutObject(&putObjectInput)

    if err != nil {
        return err
    }
    
    return nil
}
