package main

import (
	"fmt"

	"github.com/go-pg/migrations/v8"
)

const (
	createAlarmsSq = `
		CREATE SEQUENCE alarms_pk_sq
			AS INTEGER
			INCREMENT BY 1
			START WITH 1
	`
	createAlarms string = `
		CREATE TABLE alarms (
			id INTEGER PRIMARY KEY
				DEFAULT nextval('alarms_pk_sq'),
			user_id BIGINT NOT NULL,
			minutes INTEGER NOT NULL,
			created TIMESTAMP NOT NULL,
			FOREIGN KEY (user_id)
				REFERENCES users(id)
		)
	`

	deleteAlarms string = `
		DROP TABLE alarms
	`

	deleteAlarmsSq string = `
		DROP SEQUENCE alarms_pk_sq
	`
)

func init() {
	migrations.MustRegisterTx(
		func(db migrations.DB) error {
			fmt.Println("creating sequence alarms_pk_sq")
			_, err := db.Exec(createAlarmsSq)
			fmt.Println("creating alarms table")
			_, err = db.Exec(createAlarms)
			return err
		},
		func(db migrations.DB) error {
			fmt.Println("dropping alarms table")
			_, err := db.Exec(deleteUsers)
			fmt.Println("dropping sequence alarms")
			_, err = db.Exec(deleteAlarmsSq)
			return err
		},
	)
}
