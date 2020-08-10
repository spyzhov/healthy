package config

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
