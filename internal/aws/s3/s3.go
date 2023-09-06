package s3

import (
    "io"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type S3 struct {
    s3 s3iface.S3API
}

func NewS3(session *session.Session) *S3 {
    return &S3{
        s3: s3.New(session),
    }
}

type UploadInput struct {
    Bucket string
    Path string
    Image io.ReadSeeker
}

func (s *S3) Upload(input UploadInput) error {
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
