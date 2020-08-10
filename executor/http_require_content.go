package executor

import (
	"fmt"

	"github.com/spyzhov/safe"
)

type HttpArgsRequireContent struct {
	RequireMatch
	Type   HttpArgsRequireContentType `json:"type"`
	Length *RequireInteger            `json:"length"`
	JSON   []struct {
		RequireJSONPath
		RequireXPath
	} `json:"json"`
	XML  []RequireXPath `json:"xml"`
	HTML []RequireXPath `json:"html"`
}

func (a *HttpArgsRequireContent) Validate() (err error) {
	if err = a.Type.Validate(); err != nil {
		return safe.Wrap(err, "type")
	}
	if err = a.RequireMatch.Validate(); err != nil {
		return err
	}
	for _, require := range a.JSON {
		if err = require.RequireJSONPath.Validate(); err != nil {
			return safe.Wrap(err, "json")
		}
		if err = require.RequireXPath.Validate(); err != nil {
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
	return nil
}

func (a *HttpArgsRequireContent) Match(content []byte) (err error) {
	if err = a.Type.Match(content); err != nil {
		return fmt.Errorf("content: TYPE not match to be %s: %w", a.Type, err)
	}
	if a.Length != nil {
		if err = a.Length.Match("content length", len(content)); err != nil {
			return err
		}
	}
	if err = a.RequireMatch.Match("content", content); err != nil {
		return err
	}
	for _, path := range a.JSON {
		if path.RequireJSONPath.JSONPath != "" {
			if err = path.RequireJSONPath.Match("content", content); err != nil {
				return err
			}
		}
		if path.RequireXPath.XPath != "" {
			if err = path.RequireXPath.Match("content", HttpArgsRequireContentTypeJSON, content); err != nil {
				return err
			}
		}
	}
	for _, path := range a.HTML {
		if path.XPath != "" {
			if err = path.Match("content", HttpArgsRequireContentTypeHTML, content); err != nil {
				return err
			}
		}
	}
	for _, path := range a.XML {
		if path.XPath != "" {
			if err = path.Match("content", HttpArgsRequireContentTypeXML, content); err != nil {
				return err
			}
		}
	}
	return nil
}
