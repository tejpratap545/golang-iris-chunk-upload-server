package controller

import (
	"app-upload-server/service"
	"context"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"os"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
)

type initilizeRquest struct {
	FileType string `json:"fileType,omitempty" bson:"fileType" xml:"fileType" url:"file-type" form:"fileType"`
}

// Initilize Multipart Upload  godoc
// @Summary Initilize Multipart Upload
// @Description This action initiates a multipart upload and returns an upload ID. This upload ID is used to associate all of the parts in the specific multipart upload. You specify this upload ID in each of your subsequent upload part requests
// @Accept  json
// @Produce  json
// @Success 200 {object} service.MultiPartUploadResponse
// @failure 400 {string} string	"error"
// @Param bottles body initilizeRquest true "file info"
// @Router /upload/initlize [post]
func InitilizeMultiPartUpload(ctx iris.Context) {

	body := new(initilizeRquest)

	if err := ctx.ReadBody(&body); err != nil {
		ctx.StopWithText(400, "Can Not Read Body ")
	}

	key := "apk/" + uuid.New().String() + "." + body.FileType

	input := &service.CreateMultiPartRequest{
		Bucket:   os.Getenv("AWS_S3_BUCKET"),
		Key:      key,
		FileType: body.FileType,
	}

	response, err := input.CreateMultiPartUpload()
	if err != nil {
		ctx.StopWithText(400, "Can Not set  start mulipart input ")
	}

	rdb := new(service.RedisClient)
	rdb.SetClient()
	rdb.Ctx = context.Background()

	rdbData := &service.AppUploadChunk{
		StartTime: time.Now(),

		Key:         response.Key,
		Bucket:      response.Bucket,
		IsCompleted: false,
	}

	status, err := rdb.Set(*response.UploadId, rdbData)

	if err != nil {
		ctx.StopWithText(400, "Can Not set  data in redis ")
	}
	log.Println(status)

	if err != nil {
		ctx.StopWithText(400, "Can Not Create MultiPart Input ")

	}
	ctx.JSON(response)

}

// Uploads a part in a multipart upload  godoc
// @Summary Uploads a part in a multipart upload
// @Description Uploads a part in a multipart upload. In this operation, you provide part data
// @Description in your request. However, you have an option to specify your existing Amazon S3
// @Description object as a data source for the part you are uploading. To upload a part from an
// @Description existing object, you use the UploadPartCopy
// @Accept  mpfd
// @Accept  x-www-form-urlencoded
// @Produce  json
// @Success 200 {object} service.PartUploadOutput
// @Param bottles body service.PartInput true "file info"
// @Param X-Upload-Id header string true "Upload ID that identifies the multipart upload."
// @failure 400 {string} string	"error"
// @Router /upload/part [post]
func UploadPart(ctx iris.Context) {

	body := new(service.PartInput)

	file, _, _ := ctx.FormFile("chunk")
	defer file.Close()
	body.Chunk, _ = ioutil.ReadAll(file)

	part, _ := strconv.ParseInt(ctx.FormValue("partNumber"), 10, 64)
	body.PartNumber = int32(part)

	body.UploadId = ctx.GetHeader("X-Upload-Id")
	body.Bucket = os.Getenv("AWS_S3_BUCKET")

	rdb := new(service.RedisClient)
	rdb.SetClient()
	rdb.Ctx = context.Background()

	rdbData := new(service.AppUploadChunk)

	err := rdb.Get(body.UploadId, rdbData)
	if err != nil {
		ctx.StopWithText(400, "Can Not get data from redis")
	}

	body.Key = rdbData.Key

	result, err := body.UploadPart()

	if err != nil {
		ctx.StopWithText(400, "Can Not upload part")
	}

	rdbData.CompletedPart = append(rdbData.CompletedPart, *result)
	status, err := rdb.Set(body.UploadId, rdbData)

	if err != nil {
		ctx.StopWithText(400, "Can Not set  data in redis ")
	}
	log.Println(status)
	response := &service.PartUploadOutput{
		ETag:       result.ETag,
		PartNumber: result.PartNumber,
	}

	ctx.JSON(response)

}

// Completes a multipart upload  godoc
// @Summary Completes a multipart upload upload
// @Description Completes a multipart upload by assembling previously uploaded parts
// @Produce  json
// @Success 200 {object} service.CompleteMultiPartOutput
// @Param X-Upload-Id header string true "Upload ID that identifies the multipart upload."
// @failure 400 {string} string	"error"
// @Router /upload/finish [post]
func FinishMultiPartUpload(ctx iris.Context) {

	uploadId := ctx.GetHeader("X-Upload-Id")

	rdb := new(service.RedisClient)
	rdb.SetClient()
	rdb.Ctx = context.Background()

	rdbData := new(service.AppUploadChunk)

	err := rdb.Get(uploadId, rdbData)
	if err != nil {
		ctx.StopWithText(400, "Can Not get data from redis")
	}

	input := &service.CompleteMultiPartInput{
		Bucket:        rdbData.Bucket,
		Key:           rdbData.Key,
		UploadId:      uploadId,
		CompletedPart: rdbData.CompletedPart,
	}

	output, err := input.CompleteMultipartUpload()

	if err != nil {
		ctx.StopWithText(400, "Can Not Complete Multipart Output")
	}

	response := &service.CompleteMultiPartOutput{
		Location: *output.Location,
		Key:      *output.Key,
	}

	ctx.JSON(response)

}

// Upload File
// @Summary Upload File
// @Description Completes a multipart upload by assembling previously uploaded parts
// @Accept  mpfd
// @Accept  x-www-form-urlencoded
// @Produce  json
// @Success 200 {object} service.CompleteMultiPartOutput
// @Param file body string true "upload file "
// @Param location body string true "location of the file in s3 bucket default is apk"
// @failure 400 {string} string	"error"
// @Router /upload/file [post]
func FileUpload(ctx iris.Context) {
	file, fileHeader, err := ctx.FormFile("file")
	if err != nil {
		ctx.StopWithText(400, "Can Not read file ")
		return
	}
	defer file.Close()
	location := ctx.FormValueDefault("location", "apk")

	putObjectInput := new(service.PutObjectInput)
	putObjectInput.Body = file

	putObjectInput.Bucket = os.Getenv("AWS_S3_BUCKET")
	putObjectInput.Key = location + "/" + uuid.New().String() + "__" + fileHeader.Filename

	if _, err := putObjectInput.UploadFile(); err != nil {
		ctx.StopWithText(400, "Can Not Upload File ")
		return
	}

	response := &service.CompleteMultiPartOutput{
		Key: putObjectInput.Key,
	}

	ctx.JSON(response)
}
