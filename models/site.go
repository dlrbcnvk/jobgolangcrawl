package models

type Site struct {
	ID   int
	Name string
}

func NewSite(id int, name string) *Site {
	return &Site{
		ID:   id,
		Name: name,
	}
}
