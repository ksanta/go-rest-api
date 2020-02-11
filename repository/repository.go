package repository

import "github.com/ksanta/go-rest-api/domain"

// Interface for the event repository
type EventRepository interface {
	// Stores a new event, sets the ID field in the given event
	Create(event *domain.Event) error
	GetById(id int) (*domain.Event, error)
	GetAll() (*[]domain.Event, error)
	Update(event *domain.Event) error
	Delete(id int) error
}
