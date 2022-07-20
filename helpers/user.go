package helpers

import (
	"feb-cli/config"
	"fmt"
	"log"
	"strconv"

	bolt "go.etcd.io/bbolt"
)

func CheckUserAuthenticate() {
	isAuthenticate, err := IsAuthenticate()
	if err != nil {
		log.Fatal("Can not check login detail")
	}
	if !isAuthenticate {
		log.Fatal("Please login . Run feblic-cli login")
	}

}

func IsAuthenticate() (bool, error) {
	db := config.DB()
	defer db.Close()
	isAuthenticate := false
	var err error
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("user"))
		fmt.Println(b.Get([]byte("isAuthenticate")))

		isAuthenticate, err = strconv.ParseBool(string(b.Get([]byte("isAuthenticate"))))
		if err != nil {
			log.Fatal("err")
		}

		return nil

	})

	return isAuthenticate, err
}

func GetAccessToken() string {
	db := config.DB()
	defer db.Close()
	var err error
	var accessToken string
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("user"))
		accessToken = string(b.Get([]byte("accessToken")))

		return nil

	})
	if err != nil {
		log.Fatal("Please Login again")
	}
	return accessToken

}
