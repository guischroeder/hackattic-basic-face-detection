package aws

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type S3Client struct {
    S3Api s3iface.S3API
}

func NewS3Client(s3 *s3.S3) *S3Client {
    return &S3Client{
        S3Api: s3,
    }
}

type UploadInput struct {
    Bucket string
    Path string
    Image io.ReadSeeker
}

func (s *S3Client) Upload(input UploadInput) error {
    putObjectInput := s3.PutObjectInput{
        Bucket: aws.String(input.Bucket),
        Key: aws.String(input.Path),
        Body: input.Image,
    }

    _, err := s.S3Api.PutObject(&putObjectInput)

    if err != nil {
        return err
    }
    
    return nil
}
