package args

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	. "github.com/spyzhov/healthy/executor/internal"
	"github.com/spyzhov/safe"
)

type ReachText struct {
	Text *string `json:"text"`
	File *string `json:"file"`
	RN   bool    `json:"rn"`
}

func (a *ReachText) Validate() (err error) {
	if a == nil {
		return nil
	}
	if a.Text != nil && a.File != nil {
		return fmt.Errorf("only one field should be set: `text` or `file`")
	}
	if a.File != nil {
		if err = IsFileExists(*a.File); err != nil {
			return safe.Wrap(err, "file")
		}
	}
	return nil
}

func (a *ReachText) Value() (reader io.ReadCloser, err error) {
	result := new(ReaderCloserCallback)
	if a.File != nil && (*a.File) != "" {
		result.ReadCloser, err = os.Open(*a.File)
		if err != nil {
			return nil, safe.Wrap(err, "open file")
		}
	} else if a.Text != nil {
		result.ReadCloser = ioutil.NopCloser(strings.NewReader(*a.Text))
	}
	if a.RN {
		proxy := make([]byte, 0, 64)
		prev := byte(0)
		result.Callback = func(p []byte, n int) ([]byte, error) {
			proxy = proxy[:0]
			for i := 0; i < n; i++ {
				if p[i] == '\n' {
					if prev != '\r' {
						proxy = append(proxy, '\r')
					}
				}
				proxy = append(proxy, p[i])
				prev = p[i]
			}
			return proxy, nil
		}
	}
	return result, nil
}
