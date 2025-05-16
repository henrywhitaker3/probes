package probes

import (
	"errors"
	"fmt"
	"sync"
)

type Subject string

type Probes struct {
	startups  map[Subject]bool
	readies   map[Subject]bool
	healthies map[Subject]bool
	mu        *sync.Mutex

	opts ProbeOpts
}

type ProbeOpts struct {
	// The initialised value of a startup subjects (default: false)
	DefaultStartup *bool
	// The initialised value of a ready subjects (default: false)
	DefaultReady *bool
	// The initialised value of a healthy subjects (default: true)
	DefaultHealthy *bool
}

func ptr[T any](in T) *T {
	return &in
}

func New(opts ProbeOpts) *Probes {
	if opts.DefaultStartup == nil {
		opts.DefaultStartup = ptr(false)
	}
	if opts.DefaultReady == nil {
		opts.DefaultReady = ptr(false)
	}
	if opts.DefaultHealthy == nil {
		opts.DefaultHealthy = ptr(true)
	}
	return &Probes{
		startups:  map[Subject]bool{},
		readies:   map[Subject]bool{},
		healthies: map[Subject]bool{},
		opts:      opts,
		mu:        &sync.Mutex{},
	}
}

var (
	ErrUnknownStartup = errors.New("unknown startup subject")
)

func (p *Probes) Started(s Subject) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.startups[s]; !ok {
		return fmt.Errorf("%w: %s", ErrUnknownStartup, string(s))
	}
	p.startups[s] = true
	return nil
}

var (
	ErrUnkownReady = errors.New("unknown ready subject")
)

func (p *Probes) Ready(s Subject) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.readies[s]; !ok {
		return fmt.Errorf("%w: %s", ErrUnkownReady, string(s))
	}
	p.readies[s] = true
	return nil
}

func (p *Probes) NotReady(s Subject) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.readies[s]; !ok {
		return fmt.Errorf("%w: %s", ErrUnkownReady, string(s))
	}
	p.readies[s] = false
	return nil
}

var (
	ErrUnknownHealth = errors.New("unknown health subject")
)

func (p *Probes) Healthy(s Subject) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.healthies[s]; !ok {
		return fmt.Errorf("%w: %s", ErrUnknownHealth, string(s))
	}
	p.healthies[s] = true
	return nil
}

func (p *Probes) Unhealthy(s Subject) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.healthies[s]; !ok {
		return fmt.Errorf("%w: %s", ErrUnknownHealth, string(s))
	}
	p.healthies[s] = false
	return nil
}

func (p *Probes) WithStartups(s ...Subject) *Probes {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, v := range s {
		p.startups[v] = *p.opts.DefaultStartup
	}
	return p
}

func (p *Probes) WithReadies(s ...Subject) *Probes {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, v := range s {
		p.readies[v] = *p.opts.DefaultReady
	}
	return p
}

func (p *Probes) WithHealthies(s ...Subject) *Probes {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, v := range s {
		p.healthies[v] = *p.opts.DefaultHealthy
	}
	return p
}

func (p *Probes) IsStarted() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	started := true
	for _, v := range p.startups {
		if !v {
			started = false
		}
	}
	return started
}

func (p *Probes) IsHealthy() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	healthy := true
	for _, v := range p.healthies {
		if !v {
			healthy = false
		}
	}
	return healthy
}

func (p *Probes) IsReady() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	ready := true
	for _, v := range p.readies {
		if !v {
			ready = false
		}
	}
	return ready
}
