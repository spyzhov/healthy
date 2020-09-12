package args

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"path"

	. "github.com/spyzhov/healthy/executor/internal"
	"github.com/spyzhov/safe"
)

type HttpArgsForm struct {
	Values map[string]string `json:"values"`
	Files  map[string]string `json:"files"`
}

func (a *HttpArgsForm) Validate() (err error) {
	if a == nil {
		return nil
	}
	for _, filename := range a.Files {
		if err = IsFileExists(filename); err != nil {
			return safe.Wrap(err, "form")
		}
	}
	return nil
}

func (a *HttpArgsForm) SubmitForm(b *bytes.Buffer) (contentType string, err error) {
	if a == nil {
		return "", nil
	}
	if len(a.Values) != 0 || len(a.Files) != 0 {
		writer := multipart.NewWriter(b)
		defer safe.Close(writer, "multipart writer")
		for key, value := range a.Values {
			err = writer.WriteField(key, value)
			if err != nil {
				return "", fmt.Errorf("form: write field[%s] error: %w", key, err)
			}
		}
		for key, filename := range a.Files {
			w, err := writer.CreateFormFile(key, path.Base(filename))
			if err != nil {
				return "", fmt.Errorf("form: create field[%s] error: %w", key, err)
			}
			content, err := ioutil.ReadFile(filename)
			if err != nil {
				return "", fmt.Errorf("form: field[%s] read file error: %w", key, err)
			}
			_, err = w.Write(content)
			if err != nil {
				return "", fmt.Errorf("form: write field[%s] error: %w", key, err)
			}
		}
		contentType = writer.FormDataContentType()
	}
	return
}
