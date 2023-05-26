package main

import (
	"fmt"

	"github.com/go-pg/migrations/v8"
)

const (
	createUsers string = `
		CREATE TABLE users (
			id BIGINT PRIMARY KEY,
			username VARCHAR(32),
			first_name VARCHAR(20) NOT NULL,
			last_name VARCHAR(20),
			created TIMESTAMP NOT NULL
		)
	`

	deleteUsers string = `
		DROP TABLE users
	`
)

func init() {
	migrations.MustRegisterTx(
		func(db migrations.DB) error {
			fmt.Println("creating users table")
			_, err := db.Exec(createUsers)
			return err
		},
		func(db migrations.DB) error {
			fmt.Println("dropping users table")
			_, err := db.Exec(deleteUsers)
			return err
		},
	)
}
