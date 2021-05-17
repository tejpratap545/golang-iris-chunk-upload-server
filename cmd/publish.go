package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"feb-cli/helpers"
	"feb-cli/pb"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

type FilesLocation struct {
	AppFileLocation     string
	AppInfoFileLocation string
}

var filesLocation FilesLocation

// publishCmd represents the publish command
var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish App to Feblic App Store",
	Long:  `This Command is use to publish app in feblic app store. For example: `,
	Run: func(cmd *cobra.Command, args []string) {

		CheckFileInput()
		helpers.CheckUserAuthenticate()

		serverAddress := "0.0.0.0:8081"
		conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
		if err != nil {
			log.Fatal("cannot dial server: ", err)
		}

		appUploadClient := pb.NewAppUploadClient(conn)
		UploadApp(appUploadClient)
	},
}

// function to check file flag if not present in command then ask to user
func CheckFileInput() {
	var err error
	reader := bufio.NewReader(os.Stdin)
	if filesLocation.AppFileLocation == "" {
		fmt.Print("Enter app info file location : ")
		filesLocation.AppFileLocation, err = reader.ReadString('\n')
		if err != nil {
			log.Fatal("can not read app build file location")
		}
	}

	if filesLocation.AppInfoFileLocation == "" {
		fmt.Print("Enter app file location : ")
		filesLocation.AppInfoFileLocation, err = reader.ReadString('\n')
		if err != nil {
			log.Fatal("can not read app info file location")
		}
	}
}

func UploadApp(client pb.AppUploadClient) {

	//   open app location file
	app, err := os.Open(filesLocation.AppFileLocation)
	if err != nil {
		log.Fatal("Can not open app location please check file location")
	}
	fileType := path.Ext(filesLocation.AppFileLocation)
	defer app.Close()

	appStat, err := app.Stat()
	appSize := strconv.Itoa(int(appStat.Size()))

	if err != nil {

		log.Fatal("camn not read stat of the app ")
	}

	// open app information file
	jsonFile, err := os.Open(filesLocation.AppInfoFileLocation)

	if err != nil {

		log.Fatal("Can not open app information file please check file location")
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var appInfo pb.AppInfo
	json.Unmarshal([]byte(byteValue), &appInfo)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Hour)
	defer cancel()

	stream, err := client.UploadApp(ctx)

	if err != nil {
		log.Fatal("cannot upload app please try again  or check your internet information: ", err)
	}

	// send basic meta data in first request like meta data and app size
	req := &pb.UploadAppRequest{
		Data: &pb.UploadAppRequest_RequestMetaData{
			RequestMetaData: &pb.UploadAppRequestInfo{
				AccessToken: helpers.GetAccessToken(), AppSize: appSize, FileType: fileType,
			},
		},
	}

	err = stream.Send(req)
	if err != nil {
		log.Fatal("Can not send data to server ")
	}

	// send info of the app

	req = &pb.UploadAppRequest{
		Data: &pb.UploadAppRequest_AppInfo{
			AppInfo: &appInfo,
		},
	}

	err = stream.Send(req)

	if err != nil {
		log.Fatal("cannot send app info to server: ", err, stream.RecvMsg(nil))
	}

	reader := bufio.NewReader(app)
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

	rootCmd.AddCommand(publishCmd)

	publishCmd.Flags().StringVarP(&filesLocation.AppFileLocation, "sourceapp", "s", "", "App file location")
	publishCmd.Flags().StringVarP(&filesLocation.AppInfoFileLocation, "sourcefile", "p", "", "App info file location")
}
