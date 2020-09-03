package args

import (
	"github.com/spyzhov/safe"
)

type DialArgsRequire struct {
	Content *RequireContent `json:"content"`
}

func (a *DialArgsRequire) Validate() (err error) {
	if a == nil {
		return nil
	}
	if err = a.Content.Validate(); err != nil {
		return safe.Wrap(err, "require: content")
	}
	return
}

func (a *DialArgsRequire) Match(content []byte) (err error) {
	if a == nil {
		return nil
	}
	if err = a.Content.Match("content", content); err != nil {
		return err
	}
	return nil
}
