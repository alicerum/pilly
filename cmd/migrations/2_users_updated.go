package main

import (
	"fmt"

	"github.com/go-pg/migrations/v8"
)

const (
	createColumn string = `
		ALTER TABLE users
		ADD COLUMN updated TIMESTAMP
	`
	dropColumn string = `
		ALTER TABLE users
		DROP COLUMN updated
	`
)

func init() {
	migrations.MustRegisterTx(
		func(db migrations.DB) error {
			fmt.Println("creating users column updated")
			_, err := db.Exec(createColumn)
			return err
		},
		func(db migrations.DB) error {
			fmt.Println("dropping users column updated")
			_, err := db.Exec(dropColumn)
			return err
		},
	)
}
