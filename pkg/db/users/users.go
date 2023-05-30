package users

import (
	"errors"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/rs/zerolog/log"
)

type User struct {
	tableName struct{}  `pg:"users,alias:u"`
	ID        int64     `pg:"id,pk"`
	Username  string    `pg:"username"`
	FirstName string    `pg:"first_name,notnull"`
	LastName  string    `pg:"last_name"`
	Created   time.Time `pg:"created,notnull"`
	Updated   time.Time `pg:"updated"`
}

type Svc struct {
	db *pg.DB
}

func NewSvc(db *pg.DB) *Svc {
	return &Svc{
		db: db,
	}
}

func (s *Svc) GetByID(id int64) (User, error) {
	log.Debug().Int64("id", id).Msg("getting user by id")
	var u User
	err := s.db.Model(&u).Where("id = ?", id).Select()
	return u, err
}

func (s *Svc) Persist(u *User) error {
	log.Debug().Int64("id", u.ID).Msg("persisting telegram user")
	_, err := s.GetByID(u.ID)
	if err != nil && errors.Is(err, pg.ErrNoRows) {
		_, err = s.db.Model(u).Insert()
	} else {
		u.Updated = time.Now().UTC()
		_, err = s.db.Model(u).WherePK().Update()
	}
	return err
}
