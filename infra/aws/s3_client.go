package aws

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Client struct {
    s3 *s3.S3
}

func NewS3Client(session *session.Session) *S3Client {
    return &S3Client{
        s3: s3.New(session),
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

    _, err := s.s3.PutObject(&putObjectInput)

    if err != nil {
        return err
    }
    
    return nil
}
