package args

import (
	"bytes"
	"fmt"

	"github.com/antchfx/htmlquery"
	"github.com/antchfx/jsonquery"
	"github.com/antchfx/xmlquery"
)

type RequireXPath struct {
	XPath string `json:"xpath"`
	RequireMatch
}

func (a *RequireXPath) Validate() error {
	return a.RequireMatch.Validate()
}

func (a *RequireXPath) Match(name string, _type HttpArgsRequireContentType, content []byte) (err error) {
	if a.XPath != "" {
		var data []byte
		if _type.Is(HttpArgsRequireContentTypeHTML) {
			data, err = a.html(content)
		} else if _type.Is(HttpArgsRequireContentTypeXML) {
			data, err = a.xml(content)
		} else if _type.Is(HttpArgsRequireContentTypeJSON) {
			data, err = a.json(content)
		} else {
			return fmt.Errorf("%s: unsupported XPath for type %s", name, _type)
		}
		if err != nil {
			return fmt.Errorf("%s: %s: %w", name, _type, err)
		}
		return a.RequireMatch.Match(name, data)
	}
	return nil
}

func (a *RequireXPath) json(content []byte) ([]byte, error) {
	root, err := jsonquery.Parse(bytes.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("JSON parse error: %w", err)
	}
	nodes, err := jsonquery.QueryAll(root, a.XPath)
	if err != nil {
		return nil, fmt.Errorf("XPath(`%s`) query error: %w", a.XPath, err)
	}
	result := make([]byte, 0)
	for _, node := range nodes {
		result = append(result, []byte(node.InnerText())...)
	}
	return result, nil
}

func (a *RequireXPath) html(content []byte) ([]byte, error) {
	root, err := htmlquery.Parse(bytes.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("JSON parse error: %w", err)
	}
	nodes, err := htmlquery.QueryAll(root, a.XPath)
	if err != nil {
		return nil, fmt.Errorf("XPath(`%s`) query error: %w", a.XPath, err)
	}
	result := make([]byte, 0)
	for _, node := range nodes {
		result = append(result, []byte(htmlquery.InnerText(node))...)
	}
	return result, nil
}

func (a *RequireXPath) xml(content []byte) ([]byte, error) {
	root, err := xmlquery.Parse(bytes.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("JSON parse error: %w", err)
	}
	nodes, err := xmlquery.QueryAll(root, a.XPath)
	if err != nil {
		return nil, fmt.Errorf("XPath(`%s`) query error: %w", a.XPath, err)
	}
	result := make([]byte, 0)
	for _, node := range nodes {
		result = append(result, []byte(node.InnerText())...)
	}
	return result, nil
}
