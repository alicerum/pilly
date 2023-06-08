package main

import (
	"fmt"

	"github.com/go-pg/migrations/v8"
)

const (
	dailySchedulesSqCreate string = `
		CREATE SEQUENCE daily_schedules_pk_sq
			AS INTEGER
			INCREMENT BY 1
			START WITH 1
	`
	dailySchedulesSqDrop string = `
		DROP SEQUENCE daily_schedules_pk_sq
	`
	dailySchedulesCreate string = `
		CREATE TABLE daily_schedules (
			id INTEGER PRIMARY KEY DEFAULT nextval('daily_schedules_pk_sq'),
			name VARCHAR(50) NOT NULL,
			message TEXT,
			state TEXT,
			script TEXT
		)
	`
	dailySchedulesDrop string = `
		DROP TABLE daily_schedules
	`
)

func init() {
	migrations.MustRegisterTx(
		func(db migrations.DB) error {
			fmt.Println("creating table daily_schedules")
			if _, err := db.Exec(dailySchedulesSqCreate); err != nil {
				return err
			}
			_, err := db.Exec(dailySchedulesCreate)
			return err
		},
		func(db migrations.DB) error {
			fmt.Println("dropping table daily_schedules")
			if _, err := db.Exec(dailySchedulesDrop); err != nil {
				return err
			}
			_, err := db.Exec(dailySchedulesSqDrop)
			return err
		},
	)
}
