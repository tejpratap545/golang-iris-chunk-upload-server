package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"feb-cli/config"
	"feb-cli/utils"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"syscall"

	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"

	"golang.org/x/crypto/ssh/terminal"
)

type User struct {
	Username string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	AccessToken string `json:"accessToken"`
}

var user User

// loginCmd represents the login command
var loginCmd = &cobra.Command{

	Use:   "login",
	Short: "login to feblic cli",

	Run: func(cmd *cobra.Command, args []string) {
		// path, err := homedir.Dir()

		// if err != nil {
		// 	panic("can not open hone dir ")

		// }

		var err error

		reader := bufio.NewReader(os.Stdin)
		if user.Username == "" {

			fmt.Print("Enter Username: ")
			user.Password, err = reader.ReadString('\n')
			if err != nil {
				panic("can not read username")
			}

		}
		if user.Password == "" {
			fmt.Print("Enter Password: ")
			bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
			if err != nil {
				panic("can not read password")
			}

			user.Password = string(bytePassword)
		}

		// start spinner
		utils.NewSpinner()

		userJson, err := json.Marshal(user)
		if err != nil {
			panic("can not read convert json")
		}

		body := bytes.NewBuffer(userJson)

		url := fmt.Sprintf("%s/api/signin/", os.Getenv("API_URL"))

		res, err := http.Post(url, "application/json", body)

		if err != nil {
			panic("can not login please check username and password")
		}
		defer res.Body.Close()

		bytes, err := ioutil.ReadAll(res.Body)

		if err != nil {
			log.Fatalln(err)
		}

		var token Token
		json.Unmarshal(bytes, &token)

		db := config.DB()

		db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucket([]byte("user"))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			b.Put([]byte("isAuthenticate"), []byte("true"))
			b.Put([]byte("accessToken"), []byte(token.AccessToken))
			return err
		})

		utils.StopSpinner()

		fmt.Println("You Are Now sccessfully authenticate. ")

	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringVarP(&user.Username, "username", "u", "", "Enter your username")
	loginCmd.Flags().StringVarP(&user.Password, "password", "p", "", "Enter your password")
}
