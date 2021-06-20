package service

import (
	"bytes"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type CreateMultiPartRequest struct {
	Bucket   string
	Key      string
	FileType string
}

type MultiPartUploadResponse struct {
	Bucket string `json:"bucket,omitempty" xml:"bucket"`
	// Object key for which the multipart upload was initiated.
	Key string `json:"key,omitempty" xml:"key"`

	// If the bucket has a lifecycle rule configured with an action to abort incomplete
	// multipart uploads and the prefix in the lifecycle rule matches the object name
	// in the request, the response includes this header. The header indicates when the
	// initiated multipart upload becomes eligible for an abort operation. For more
	// information, see  Aborting Incomplete Multipart Uploads Using a Bucket Lifecycle
	// Policy
	// (https://docs.aws.amazon.com/AmazonS3/latest/dev/mpuoverview.html#mpu-abort-incomplete-mpu-lifecycle-config).
	// The response also includes the x-amz-abort-rule-id header that provides the ID
	// of the lifecycle configuration rule that defines this action.
	AbortDate *time.Time `json:"abortDate,omitempty" xml:"abortDate"`

	// ID for the initiated multipart upload.
	UploadId *string `json:"uploadId,omitempty" xml:"uploadId"`
}

// input data for upload chunk
type PartInput struct {
	Chunk      []byte `json:"chunk,omitempty" xml:"chunk" form:"chunk" url:"chunk"`
	Bucket     string `json:"-" xml:"-" form:"-" url:"-"`
	PartNumber int32  `json:"partNumber,omitempty" xml:"partNumber" form:"partNumber" url:"partNumber"`
	UploadId   string `json:"-" xml:"-" form:"-" url:"-"`
	Key        string `json:"-" xml:"-" form:"-" url:"-"`
}

// output data after upload chunk
type PartUploadOutput struct {

	// Entity tag returned when the part was uploaded.
	ETag *string

	// Part number that identifies the part. This is a positive integer between 1 and
	// 10,000.
	PartNumber int32
}

func (input *PartInput) GetBody() *bytes.Reader {
	return bytes.NewReader(input.Chunk)
}
func (input *PartInput) Length() int {
	return len(input.Chunk)
}

// data for complete multipart upload
type CompleteMultiPartInput struct {
	Bucket        string
	Key           string
	UploadId      string
	CompletedPart []types.CompletedPart
}

type CompleteMultiPartOutput struct {
	Location string `json:"location,omitempty" xml:"location"`
	Key      string `json:"key,omitempty"`
}

// data store in redis db
type AppUploadChunk struct {
	Key           string
	Bucket        string
	IsCompleted   bool
	StartTime     time.Time
	AbortTime     time.Time
	CompletedPart []types.CompletedPart
}
