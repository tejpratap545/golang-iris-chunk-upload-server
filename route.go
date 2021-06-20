package main

import (
	"app-upload-server/controller"

	"github.com/kataras/iris/v12"
)

func Route(route iris.Party) {

	route.HandleDir("/docs", iris.Dir("docs"))

	upoad := route.Party("/upload")
	upoad.Post("/initilize", controller.InitilizeMultiPartUpload)
	upoad.Post("/part", controller.UploadPart)

	upoad.Post("/finish", controller.FinishMultiPartUpload)
	upoad.Post("/file", controller.FileUpload)

}
