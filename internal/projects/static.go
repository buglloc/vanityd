package projects

import (
	"fmt"
	"slices"

	"github.com/goccy/go-yaml"
)

var _ Provider = (*StaticProjects)(nil)

type StaticProjects struct {
	projects  []Project
	slugIndex map[string]int
}

func NewStaticProjects(rawData []byte) (*StaticProjects, error) {
	out := &StaticProjects{
		slugIndex: make(map[string]int),
	}

	if err := yaml.Unmarshal(rawData, &out.projects); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	for i, prj := range out.projects {
		for _, slug := range prj.Slug {
			out.slugIndex[slug] = i
		}
	}

	return out, nil
}

func (p *StaticProjects) List() []Project {
	return slices.Clone(p.projects)
}

func (p *StaticProjects) Get(slug string) (Project, error) {
	idx, ok := p.slugIndex[slug]
	if !ok {
		return Project{}, ErrNotFound
	}

	return p.projects[idx], nil
}
