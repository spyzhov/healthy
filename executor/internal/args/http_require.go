package args

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spyzhov/safe"
)

type HttArgsRequire struct {
	Status  *HttpArgsRequireStatus  `json:"status"`
	Content *HttpArgsRequireContent `json:"content"`
	Header  *HttpArgsRequireHeader  `json:"header"`
}

func (a *HttArgsRequire) Validate() (err error) {
	if a.Status != nil {
		if err = a.Status.Validate(); err != nil {
			return safe.Wrap(err, "require: status")
		}
	}
	if a.Content != nil {
		if err = a.Content.Validate(); err != nil {
			return safe.Wrap(err, "require: content")
		}
	}
	if a.Header != nil {
		if err = a.Header.Validate(); err != nil {
			return safe.Wrap(err, "require: header")
		}
	}
	return
}

func (a *HttArgsRequire) Match(response *http.Response) (err error) {
	if a.Status != nil {
		err = a.Status.Match(response.StatusCode)
		if err != nil {
			return err
		}
	}
	if a.Header != nil {
		err = a.Header.Match(response.Header)
		if err != nil {
			return err
		}
	}
	if a.Content != nil {
		var content []byte
		content, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("read: %w", err)
		}
		err = a.Content.Match(content)
		if err != nil {
			return err
		}
	}
	return nil
}
