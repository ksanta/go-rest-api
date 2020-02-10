package domain

type Event struct {
	ID          int    `json:"ID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}
