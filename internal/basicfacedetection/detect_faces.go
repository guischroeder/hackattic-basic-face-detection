package basicfacedetection

import (
    "io"
    "net/http"
    "os"
    "strings"

    "hackattic-basic-face-detection/internal/hackattic"
    "hackattic-basic-face-detection/internal/aws/s3"
    "hackattic-basic-face-detection/internal/aws/rekognition"
)

const PATH = "media/faces.jpg"

type DetectFaces struct {
    s3 *s3.S3
    rekognition *rekognition.Rekognition
    HttpClient http.Client
}
func NewDetectFaces(s3 *s3.S3, rekognition *rekognition.Rekognition) *DetectFaces {
    return &DetectFaces{s3: s3, rekognition: rekognition}
}

func (d *DetectFaces) Perform(problem hackattic.Problem) (*rekognition.DetectFacesOutput, error) {
    bucket := os.Getenv("S3_BUCKET_NAME")

    err := d.uploadImageFromUrlToS3(problem.ImageUrl, bucket, PATH)
    if err != nil {
        return nil, err
    }

    detectedFaces, err := d.rekognition.DetectFaces(bucket, PATH)
    if err != nil {
        return nil, err
    }
    
    return detectedFaces, nil
}

func (d *DetectFaces) uploadImageFromUrlToS3(
    imageUrl string, bucket string, path string) error {
	request, err := http.NewRequest(http.MethodGet, imageUrl, nil)

    response, err := d.HttpClient.Do(request)

    if err != nil {
        return err
    }
    defer response.Body.Close()

    data, _ := io.ReadAll(response.Body)
    reader := strings.NewReader(string(data))

    err = d.s3.Upload(s3.UploadInput{
        Bucket: bucket,
        Path: path,
        Image: io.ReadSeeker(reader),
    })
    if err != nil {
        return err
    }

    return nil
}
