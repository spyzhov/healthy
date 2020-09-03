package args

import (
	"fmt"

	"github.com/spyzhov/safe"
)

type RequireContent struct {
	Type   HttpArgsRequireContentType `json:"type"`
	Length *RequireNumeric            `json:"length"`
	JSON   []RequireJSON              `json:"json"`
	XML    []RequireXPath             `json:"xml"`
	HTML   []RequireXPath             `json:"html"`
	Text   []RequireMatch             `json:"text"`
}

func (a *RequireContent) Validate() (err error) {
	if a == nil {
		return nil
	}
	if err = a.Type.Validate(); err != nil {
		return safe.Wrap(err, "type")
	}
	if a.Length != nil {
		if err = a.Length.Validate(); err != nil {
			return safe.Wrap(err, "length")
		}
	}
	for _, require := range a.JSON {
		if err = require.Validate(); err != nil {
			return safe.Wrap(err, "json")
		}
	}
	for _, require := range a.HTML {
		if err = require.Validate(); err != nil {
			return safe.Wrap(err, "html")
		}
	}
	for _, require := range a.XML {
		if err = require.Validate(); err != nil {
			return safe.Wrap(err, "xml")
		}
	}
	for _, require := range a.Text {
		if err = require.Validate(); err != nil {
			return safe.Wrap(err, "text")
		}
	}
	return nil
}

func (a *RequireContent) Match(name string, content []byte) (err error) {
	if a == nil {
		return nil
	}
	if err = a.Type.Match(content); err != nil {
		return fmt.Errorf("%s: TYPE not match to be %s: %w", name, a.Type, err)
	}
	if a.Length != nil {
		if err = a.Length.Match(name+": length", float64(len(content))); err != nil {
			return err
		}
	}
	for _, path := range a.JSON {
		if err = path.Match(name, content); err != nil {
			return err
		}
	}
	for _, path := range a.HTML {
		if path.XPath != "" {
			if err = path.Match(name, HttpArgsRequireContentTypeHTML, content); err != nil {
				return err
			}
		}
	}
	for _, path := range a.XML {
		if path.XPath != "" {
			if err = path.Match(name, HttpArgsRequireContentTypeXML, content); err != nil {
				return err
			}
		}
	}
	for _, require := range a.Text {
		if err = require.Match(name+": text", content); err != nil {
			return err
		}
	}
	return nil
}
