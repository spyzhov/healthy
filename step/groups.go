package step

import "sync"

type Groups struct {
	mu    sync.Mutex
	steps map[string]*Steps
	names []string
}

func NewGroups() *Groups {
	return &Groups{
		steps: make(map[string]*Steps),
		names: make([]string, 0),
	}
}

func (g *Groups) Names() []string {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.names
}

func (g *Groups) Get(name string) *Steps {
	g.mu.Lock()
	defer g.mu.Unlock()
	if step, ok := g.steps[name]; ok {
		return step
	}
	g.names = append(g.names, name)
	g.steps[name] = NewSteps()
	return g.steps[name]
}
