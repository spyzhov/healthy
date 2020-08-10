package step

import "sync"

type Steps struct {
	mu    sync.Mutex
	steps []*Step
	index map[string]int
	names []string
}

func NewSteps() *Steps {
	return &Steps{
		steps: make([]*Step, 0),
		index: make(map[string]int),
		names: make([]string, 0),
	}
}

func (s *Steps) Names() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.names
}

func (s *Steps) Add(name string, test Function) *Steps {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.names = append(s.names, name)
	s.index[name] = len(s.steps)
	s.steps = append(s.steps, &Step{
		Name: name,
		Func: test,
	})
	return s
}

func (s *Steps) Get(name string) *Step {
	s.mu.Lock()
	defer s.mu.Unlock()
	index, ok := s.index[name]
	if ok {
		return s.steps[index]
	}
	return nil
}
