package services

import (
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/service/rekognition"

	"hackattic-basic-face-detection/infra/aws"
	"hackattic-basic-face-detection/infra/hackattic"
	"hackattic-basic-face-detection/infra/http_client"
)

const PATH = "media/faces.jpg"

type DetectFaces struct {
    s3Client *aws.S3Client
    rekognitionClient *aws.RekognitionClient
}
func NewDetectFaces(s3Client *aws.S3Client, rekognitionClient *aws.RekognitionClient) *DetectFaces {
    return &DetectFaces{s3Client: s3Client, rekognitionClient: rekognitionClient}
}

func (d *DetectFaces) Perform(problem hackattic.Problem) (*rekognition.DetectFacesOutput, error) {
    bucket := os.Getenv("S3_BUCKET_NAME")

    err := d.uploadImageFromUrlToS3(problem.ImageUrl, bucket, PATH)
    if err != nil {
        return nil, err
    }

    detectedFaces, err := d.rekognitionClient.DetectFaces(aws.S3Object{
        Bucket: bucket,
        Name: PATH,
    })
    if err != nil {
        return nil, err
    }
    
    return detectedFaces, nil
}

func (d *DetectFaces) uploadImageFromUrlToS3(
    imageUrl string, bucket string, path string) error {
    resp, err := httpclient.Client.Get(imageUrl)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    data, _ := io.ReadAll(resp.Body)
    reader := strings.NewReader(string(data))

    err = d.s3Client.Upload(aws.UploadInput{
        Bucket: bucket,
        Path: path,
        Image: io.ReadSeeker(reader),
    })
    if err != nil {
        return err
    }

    return nil
}
