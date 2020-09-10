package config

import (
	"encoding/json"
	"fmt"
)

type Config struct {
	Version  int      `json:"version"`
	Name     string   `json:"name"`
	Groups   []Group  `json:"groups"`
	Frontend Frontend `json:"frontend"`
}

type Group struct {
	Name     string `json:"name"`
	Validate []Step `json:"validate"`
}

type Step struct {
	Name string        `json:"name"`
	Type string        `json:"type"`
	Vars Variables     `json:"variables"`
	Args []interface{} `json:"args"`
}

type Frontend struct {
	Script FrontendPart `json:"script"`
	Style  FrontendPart `json:"style"`
}

type FrontendPart struct {
	Content string   `json:"content"`
	Files   []string `json:"files"`
}

type Variable struct {
	Name   string      `json:"name"`
	Value  interface{} `json:"value"`
	Masked bool        `json:"masked"`
}

type Variables []*Variable

type step struct {
	Name string        `json:"name"`
	Type string        `json:"type"`
	Vars Variables     `json:"variables"`
	Args []interface{} `json:"args"`
}

func (s *Step) UnmarshalJSON(data []byte) (err error) {
	raw := new(step)
	err = raw.unmarshal(data, make(map[string]string))
	if err != nil {
		return err
	}
	raw.SetStep(s)
	return nil
}

func (s *step) unmarshal(data []byte, vars map[string]string) (err error) {
	err = json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	for _, variable := range s.Vars {
		variable.Value = expand(variable.value(), vars)
		vars[variable.Name] = variable.value()
	}
	switch s.Type {
	case "repeat":
		if len(s.Args) > 1 {
			s.Args[1], err = stepDereference(s.Args[1], vars)
			if err != nil {
				return err
			}
		}
	}
	for i, arg := range s.Args {
		s.Args[i] = dereference(arg, func(str string) string {
			return expand(str, vars)
		})
	}
	return nil
}

func (s *step) SetStep(st *Step) {
	st.Name = s.Name
	st.Type = s.Type
	st.Vars = s.Vars
	st.Args = s.Args
}

func (v Variables) Masked() []string {
	result := make([]string, 0)
	for _, variable := range v {
		if variable.Masked {
			result = append(result, variable.value())
		}
	}
	return result
}

func (v *Variable) value() string {
	return fmt.Sprint(v.Value)
}
