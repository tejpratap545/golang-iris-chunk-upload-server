package main

import (
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
)

func main() {
	app := newApp()

	crs := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*", "Content-Range", "X-Upload-Id"},
	})

	app.UseRouter(crs)

	app.Listen(":8081")
}

func newApp() *iris.Application {

	app := iris.New()

	app.PartyFunc("/", Route)

	return app
}
