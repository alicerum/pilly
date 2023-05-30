package db

import (
	"github.com/alicerum/pilly/pkg/db/daily/alarms"
	"github.com/alicerum/pilly/pkg/db/users"
	"github.com/go-pg/pg/v10"
)

type Svc struct {
	usersSvc  *users.Svc
	alarmsSvc *alarms.Svc
}

func New(db *pg.DB) *Svc {
	return &Svc{
		usersSvc:  users.NewSvc(db),
		alarmsSvc: alarms.NewSvc(db),
	}
}

func (s *Svc) Users() *users.Svc {
	return s.usersSvc
}

func (s *Svc) Alarms() *alarms.Svc {
	return s.alarmsSvc
}
