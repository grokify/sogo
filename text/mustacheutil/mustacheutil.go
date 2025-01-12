package mustacheutil

import (
	"bytes"
	"fmt"

	"github.com/cbroglie/mustache"
)

type MustacheSet struct {
	Filenames map[string]string
	Templates map[string]*mustache.Template
}

func (ms *MustacheSet) ReadTemplates() error {
	for key, filename := range ms.Filenames {
		if filename == "" {
			continue
		}
		if tmpl, err := mustache.ParseFile(filename); err != nil {
			return err
		} else {
			if ms.Templates == nil {
				ms.Templates = map[string]*mustache.Template{}
			}
			ms.Templates[key] = tmpl
		}
	}
	return nil
}

func (ms *MustacheSet) RenderTemplate(key string, data map[string]string) (*bytes.Buffer, error) {
	tmpl, ok := ms.Templates[key]
	if !ok {
		return nil, fmt.Errorf("template key not present for key (%s)", key)
	} else if tmpl == nil {
		return nil, fmt.Errorf("template is nil for key (%s)", key)
	} else {
		var buf bytes.Buffer
		if err := tmpl.FRender(&buf, data); err != nil {
			return nil, err
		} else {
			return &buf, nil
		}
	}
}
