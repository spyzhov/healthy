package args

import (
	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/spyzhov/safe"
)

type AppConsulArgsRequireService struct {
	Tags    []string `json:"tags"`
	Exist   *Bool    `json:"exist"`
	Healthy *Bool    `json:"healthy"`
}

func (a *AppConsulArgsRequireService) Validate() (err error) {
	if a == nil {
		return nil
	}
	if err = a.Exist.Validate(); err != nil {
		return safe.Wrap(err, "exist")
	}
	return
}

func (a *AppConsulArgsRequireService) Match(out []*api.ServiceEntry) (err error) {
	if a == nil {
		return nil
	}
	if err = a.Exist.Match(len(out) > 0, "not exist", "exist"); err != nil {
		return safe.Wrap(err, "exist")
	}
	for _, entry := range out {
		status := entry.Checks.AggregatedStatus()
		message := fmt.Sprintf("wrong status: %s", status)
		if err := a.Healthy.Match(status == api.HealthPassing, message, message); err != nil {
			return safe.Wrap(err, "healthy")
		}
	}
	return
}
