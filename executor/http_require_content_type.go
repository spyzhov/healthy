package executor

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"

	"golang.org/x/net/html"
	"gopkg.in/yaml.v2"
)

type HttpArgsRequireContentType string

const (
	HttpArgsRequireContentTypeJSON HttpArgsRequireContentType = "JSON"
	HttpArgsRequireContentTypeXML  HttpArgsRequireContentType = "XML"
	HttpArgsRequireContentTypeHTML HttpArgsRequireContentType = "HTML"
	HttpArgsRequireContentTypeYAML HttpArgsRequireContentType = "YAML"
)

func (v HttpArgsRequireContentType) Validate() error {
	value := HttpArgsRequireContentType(strings.ToUpper(string(v)))
	switch value {
	case "",
		HttpArgsRequireContentTypeHTML,
		HttpArgsRequireContentTypeJSON,
		HttpArgsRequireContentTypeXML,
		HttpArgsRequireContentTypeYAML:
		return nil
	default:
		return fmt.Errorf("type should be one of: %v", []HttpArgsRequireContentType{
			HttpArgsRequireContentTypeHTML,
			HttpArgsRequireContentTypeJSON,
			HttpArgsRequireContentTypeXML,
			HttpArgsRequireContentTypeYAML,
		})
	}
}

func (v HttpArgsRequireContentType) Match(content []byte) error {
	value := HttpArgsRequireContentType(strings.ToUpper(string(v)))
	switch value {
	case "":
		return nil
	case HttpArgsRequireContentTypeHTML:
		_, err := html.Parse(bytes.NewReader(content))
		return err
	case HttpArgsRequireContentTypeJSON:
		return json.Unmarshal(content, new(interface{}))
	case HttpArgsRequireContentTypeXML:
		return xml.Unmarshal(content, new(interface{}))
	case HttpArgsRequireContentTypeYAML:
		return yaml.Unmarshal(content, new(interface{}))
	default:
		return fmt.Errorf("type should be one of: %v", []HttpArgsRequireContentType{
			HttpArgsRequireContentTypeHTML,
			HttpArgsRequireContentTypeJSON,
			HttpArgsRequireContentTypeXML,
			HttpArgsRequireContentTypeYAML,
		})
	}
}

func (v HttpArgsRequireContentType) Is(value HttpArgsRequireContentType) bool {
	return value == HttpArgsRequireContentType(strings.ToUpper(string(v)))
}
