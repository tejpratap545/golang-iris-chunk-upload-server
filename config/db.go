package config

import bolt "go.etcd.io/bbolt"

func DB() *bolt.DB {

	db, err := bolt.Open(".data", 0666, nil)
	if err != nil {
		panic("can not open hone dir ")
	}

	defer db.Close()
	return db
}
