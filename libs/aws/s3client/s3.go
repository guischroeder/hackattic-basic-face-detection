package s3client

import (
		"github.com/aws/aws-sdk-go/aws/session"
		"github.com/aws/aws-sdk-go/service/s3"
)

func S3(sess *session.Session) (*s3.S3) {
    return s3.New(sess)
}

