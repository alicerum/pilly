package alarms

import (
	"time"

	"github.com/alicerum/pilly/pkg/db/users"
	"github.com/go-pg/pg/v10"
)

type Alarm struct {
	tableName struct{}  `pg:"alarms"`
	ID        int       `pg:"id,pk"`
	UserID    int64     `pg:"user_id"`
	Minutes   int       `pg:"minutes"`
	Created   time.Time `pg:"created"`

	User *users.User `pg:"rel:has-one"`
}

type Svc struct {
	db *pg.DB
}

func NewSvc(db *pg.DB) *Svc {
	return &Svc{
		db: db,
	}
}

func (s *Svc) GetByUser(id int64) ([]Alarm, error) {
	var alarms []Alarm
	err := s.db.Model(&alarms).
		Relation("User").
		Where("user_id = ?", id).
		Select()

	if err != nil {
		return nil, err
	}
	return alarms, nil
}

func (s *Svc) GetByID(id int) (*Alarm, error) {
	var alarm Alarm
	err := s.db.Model(&alarm).
		Relation("User").
		Where("Alarm.id = ?", id).
		Select()
	if err != nil {
		return nil, err
	}
	return &alarm, nil
}

func (s *Svc) Persist(alarm *Alarm) error {
	alarm.Created = time.Now().UTC()
	_, err := s.db.Model(alarm).Insert()
	return err
}
