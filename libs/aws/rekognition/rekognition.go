package rekognition

import (
		"github.com/aws/aws-sdk-go/aws/session"
		"github.com/aws/aws-sdk-go/service/rekognition"
)

func Rekognition(sess *session.Session) (*rekognition.Rekognition) {
    return rekognition.New(sess)
}

