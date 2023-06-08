package schedule

import "github.com/go-pg/pg/v10"

type Schedule struct {
	tableName struct{} `pg:"daily_schedules"`
	ID        int      `pg:"id,pk"`
	Name      string   `pg:"name,notnull"`
	Message   string   `pg:"message"`
	State     string   `pg:"state"`
	Script    string   `pg:"script"`
}

type Svc struct {
	db *pg.DB
}

func NewSvc(db *pg.DB) *Svc {
	return &Svc{
		db: db,
	}
}

func (s *Svc) GetByID(id int) (*Schedule, error) {
	var sch Schedule
	err := s.db.Model(&sch).Where("id = ?").Select()
	if err != nil {
		return nil, err
	}
	return &sch, nil
}

func (s *Svc) Persist(schedule *Schedule) error {
	var err error
	if schedule.ID == 0 {
		_, err = s.db.Model(schedule).Insert()
	} else {
		_, err = s.db.Model(schedule).WherePK().Update()
	}
	return err
}
