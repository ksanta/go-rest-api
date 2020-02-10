package repository

import (
	"database/sql"
	"github.com/ksanta/go-rest-api/domain"
	"log"
)

type eventRepoImpl struct {
	db *sql.DB
}

// Creates an object that represents the EventRepository interface
func NewPostgresEventRepo() EventRepository {
	connStr := "user=postgres password=password dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return eventRepoImpl{db}
}

func (e eventRepoImpl) Create(event *domain.Event) error {
	panic("implement me")
}

func (e eventRepoImpl) GetById(id int) (*domain.Event, error) {
	panic("implement me")
}

func (e eventRepoImpl) GetAll() ([]*domain.Event, error) {
	panic("implement me")
}

func (e eventRepoImpl) Update(event *domain.Event) error {
	panic("implement me")
}

func (e eventRepoImpl) Delete(id int) error {
	panic("implement me")
}
