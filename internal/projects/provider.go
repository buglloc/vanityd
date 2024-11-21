package projects

import "errors"

var ErrNotFound = errors.New("project not found")

type Provider interface {
	List() []Project
	Get(slug string) (Project, error)
}

type Project struct {
	Name        string
	Description string
	Slug        []string
	URL         string
}
