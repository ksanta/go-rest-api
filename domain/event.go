package domain

import "errors"

type Event struct {
	ID          int    `json:"ID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}

// Returns nil if there are no validation failures
func (e Event) Validate() error {
	if e.Title == "" {
		return errors.New("title cannot be empty")
	}
	if e.Description == "" {
		return errors.New("description cannot be empty")
	}
	return nil
}
