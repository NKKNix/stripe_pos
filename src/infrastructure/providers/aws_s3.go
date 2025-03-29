package providers

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type AwsS3Upload struct {
}

type IAwsS3Upload interface {
	HashString(id string) string
	CreateKeyNameImage(fileName string, typeFile string) (string, string)
	UploadS3FromString(fileName []byte, keyName string, contentType string) (string, error)
}

func NewAwsS3() IAwsS3Upload {
	return &AwsS3Upload{}
}

func (s3 *AwsS3Upload) HashString(id string) string {
	byteID := []byte(id)
	hashObject := sha256.Sum256(byteID)
	hexDig := hex.EncodeToString(hashObject[:])
	return hexDig
}

func (s *AwsS3Upload) CreateKeyNameImage(fileName string, typeFile string) (string, string) {
	keyName := fmt.Sprintf("images/%v.%v", strings.ToUpper(fileName), typeFile)
	contentType := fmt.Sprintf("image/%v", typeFile)
	return keyName, contentType
}

func (s *AwsS3Upload) UploadS3FromString(fileName []byte, keyName string, contentType string) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), ""),
	})
	if err != nil {
		return "", err
	}

	svc := s3.New(sess)
	var BUCKET = os.Getenv("AWS_BUCKET_NAME")

	_, err = svc.PutObject(&s3.PutObjectInput{
		Body:        bytes.NewReader(fileName),
		Bucket:      aws.String(BUCKET),
		Key:         aws.String(keyName),
		ContentType: aws.String(contentType),
		Metadata:    map[string]*string{"Content-Disposition": aws.String("attachment")},
		ACL:         aws.String("public-read"),
	})
	if err != nil {
		return "", err
	}

	fullURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", BUCKET, os.Getenv("AWS_REGION"), keyName)

	return fullURL, nil
}
