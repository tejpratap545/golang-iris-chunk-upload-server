package cmd

import (
	"bufio"
	"context"
	"feb-cli/pb"
	"io"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var appFileLocation string
var appInfoFileLocation string

// publishCmd represents the publish command
var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish App to Feblic App Store",
	Long:  `This Command is use to publish app in feblic app store. For example: `,
	Run: func(cmd *cobra.Command, args []string) {

		serverAddress := "0.0.0.0:8081"
		conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
		if err != nil {
			log.Fatal("cannot dial server: ", err)
		}

		appUploadClient := pb.NewAppUploadClient(conn)
		UploadApp(appUploadClient)
	},
}

func UploadApp(client pb.AppUploadClient) {
	filePath := "tmp/laptop.jpg"
	file, err := os.Open(filePath)
	// fileInfo, _ := file.Stat()
	// fileType := path.Ext(filePath)

	if err != nil {
		log.Fatal("cannot open image file: ", err)
	}
	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Hour)
	defer cancel()

	stream, err := client.UploadApp(ctx)

	if err != nil {
		log.Fatal("cannot upload image: ", err)
	}

	// err = stream.Send(req)

	if err != nil {
		log.Fatal("cannot send image info to server: ", err, stream.RecvMsg(nil))
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, 5*1024*1024)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("cannot read chunk to buffer: ", err)
		}

		req := &pb.UploadAppRequest{
			Data: &pb.UploadAppRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		err = stream.Send(req)
		if err != nil {
			log.Fatal("cannot send chunk to server: ", err, stream.RecvMsg(nil))
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("cannot receive response: ", err)
	}

	log.Printf("app uploaded with id: %s, size: %f", res.GetUrl(), res.GetSize())
}

func init() {

	publishCmd.Flags().StringVarP(&appFileLocation, "sourceapp", "s", "", "Source directory to read from")
	publishCmd.Flags().StringVarP(&appInfoFileLocation, "sourcefile", "p", "", "Source directory to read from")
}
