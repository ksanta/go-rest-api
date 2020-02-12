package repository

import (
	"database/sql"
	"github.com/ksanta/go-rest-api/domain"
	_ "github.com/lib/pq"
)

type eventRepoImpl struct {
	db *sql.DB
}

// Creates an object that represents the EventRepository interface
func NewPostgresEventRepo() (EventRepository, error) {
	eventRepo := &eventRepoImpl{}
	err := eventRepo.initialise()
	if err != nil {
		return nil, err
	}
	return eventRepo, nil
}

func (e *eventRepoImpl) initialise() (err error) {
	connStr := "user=postgres password=password dbname=postgres sslmode=disable"
	e.db, err = sql.Open("postgres", connStr)
	return err
}

// Stores the given event into Postgres, and sets the ID field
func (e eventRepoImpl) Create(event *domain.Event) error {
	result, err := e.db.Exec("insert into events (title, description) VALUES ($1, $2)", event.Title, event.Description)
	if err != nil {
		return err
	}

	// todo: LastInsertId() is not supported by the Postgres driver
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return err
	}

	event.ID = int(lastInsertId)
	return nil
}

// Returns nil, nil if no rows are found, else returns the event, or an error
func (e eventRepoImpl) GetById(id int) (*domain.Event, error) {
	row := e.db.QueryRow("SELECT id, title, description FROM events where id = $1", id)
	event := domain.Event{}

	err := row.Scan(&event.ID, &event.Title, &event.Description)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return &event, nil
	}
}

func (e eventRepoImpl) GetAll() (*[]domain.Event, error) {
	rows, err := e.db.Query("SELECT id, title, description FROM events")
	if err != nil {
		return nil, err
	}

	events := make([]domain.Event, 0)

	for rows.Next() {
		event := domain.Event{}

		err = rows.Scan(&event.ID, &event.Title, &event.Description)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return &events, nil
}

func (e eventRepoImpl) Update(event *domain.Event) error {
	_, err := e.db.Exec("update events set title = $2, description = $3 where id = $1", event.ID, event.Title, event.Description)
	return err
}

func (e eventRepoImpl) Delete(id int) error {
	_, err := e.db.Exec("delete from events where id = $1", id)
	return err
}
