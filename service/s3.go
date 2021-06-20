package service

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func NewS3Client() *s3.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_S3_REGION")),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: os.Getenv("AWS_ACCESS_KEY_ID"), SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"), SessionToken: "",
			},
		}))
	if err != nil {
		log.Fatal(err)
	}
	client := s3.NewFromConfig(cfg)

	return client

}

func (request *CreateMultiPartRequest) CreateMultiPartUpload() (*MultiPartUploadResponse, error) {
	s3Client := NewS3Client()

	result, err := s3Client.CreateMultipartUpload(context.TODO(), &s3.CreateMultipartUploadInput{
		Bucket:      &request.Bucket,
		Key:         &request.Key,
		ContentType: &request.FileType,
	})

	if err != nil {
		return nil, err
	}

	return &MultiPartUploadResponse{
		Bucket:    *result.Bucket,
		Key:       *result.Key,
		AbortDate: result.AbortDate,
		UploadId:  result.UploadId,
	}, err

}

func (request *PartInput) UploadPart() (*types.CompletedPart, error) {

	s3Client := NewS3Client()

	partInpurt := &s3.UploadPartInput{
		Body:          request.GetBody(),
		Bucket:        &request.Bucket,
		PartNumber:    request.PartNumber,
		UploadId:      &request.UploadId,
		ContentLength: int64(request.Length()),
		Key:           &request.Key,
	}

	result, err := s3Client.UploadPart(context.TODO(), partInpurt)
	if err != nil {
		return nil, err
	}

	return &types.CompletedPart{
		ETag:       result.ETag,
		PartNumber: request.PartNumber,
	}, err

}

func (input *CompleteMultiPartInput) CompleteMultipartUpload() (*s3.CompleteMultipartUploadOutput, error) {

	s3Client := NewS3Client()

	completeInput := &s3.CompleteMultipartUploadInput{
		Bucket:          &input.Bucket,
		Key:             &input.Key,
		UploadId:        &input.UploadId,
		MultipartUpload: &types.CompletedMultipartUpload{Parts: input.CompletedPart},
	}

	return s3Client.CompleteMultipartUpload(context.TODO(), completeInput)
}

type PutObjectInput struct {
	Key    string
	Bucket string
	Body   io.Reader
}

func (input *PutObjectInput) UploadFile() (*s3.PutObjectOutput, error) {

	s3Client := NewS3Client()
	return s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &input.Bucket,
		Key:    &input.Key,
		Body:   input.Body,
	})
}
