package cmd

import (
	"feb-cli/config"
	"fmt"

	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout to feblic cli ",

	Run: func(cmd *cobra.Command, args []string) {
		db := config.DB()

		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("user"))
			err := b.Put([]byte("isAuthenticate"), []byte("false"))
			b.Put([]byte("accessToken"), []byte(""))
			return err
		})
		fmt.Println("You are successfully logout")

	},
}
