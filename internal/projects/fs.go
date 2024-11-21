package projects

import (
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"

	"github.com/rs/zerolog/log"
)

var _ Provider = (*FSProjects)(nil)

type FSProjects struct {
	projects   *StaticProjects
	sourcePath string
	epoch      int64
	reloadMu   sync.Mutex
}

func NewFSProjects(sourcePath string) (*FSProjects, error) {
	prj := &FSProjects{
		sourcePath: sourcePath,
	}

	return prj, prj.Reload()
}

func (p *FSProjects) Reload() error {
	p.reloadMu.Lock()
	defer p.reloadMu.Unlock()

	if !p.NeedReload() {
		return nil
	}

	log.Info().Msg("reloading projects")
	f, err := os.Open(p.sourcePath)
	if err != nil {
		return fmt.Errorf("open projects file: %w", err)
	}
	defer func() { _ = f.Close() }()

	stat, err := f.Stat()
	if err != nil {
		return fmt.Errorf("stat projects file: %w", err)
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("read projects file: %w", err)
	}

	prj, err := NewStaticProjects(data)
	if err != nil {
		return fmt.Errorf("parse projects file: %w", err)
	}

	p.projects = prj
	atomic.StoreInt64(&p.epoch, stat.ModTime().Unix())
	return nil
}

func (p *FSProjects) List() []Project {
	p.TryReload()
	return p.projects.List()
}

func (p *FSProjects) Get(slug string) (Project, error) {
	p.TryReload()
	return p.projects.Get(slug)
}

func (p *FSProjects) TryReload() {
	if !p.NeedReload() {
		return
	}

	if err := p.Reload(); err != nil {
		log.Error().Err(err).Msg("reload projects")
	}
}

func (p *FSProjects) NeedReload() bool {
	st, err := os.Stat(p.sourcePath)
	if err != nil {
		return false
	}

	curEpoch := atomic.LoadInt64(&p.epoch)
	newEpoch := st.ModTime().Unix()
	return newEpoch > curEpoch
}
