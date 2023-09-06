package aws

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws/session"
)

func NewSession(key string, secret string, region string) (*session.Session, error) {
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String(region),
        Credentials: credentials.NewStaticCredentials(key, secret, "",),
    })

    if err != nil { 
        return nil, err 
    }

    return sess, nil
}
