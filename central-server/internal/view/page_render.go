package view

import (
	"fmt"
	"html/template"
	"io"
)

type Render interface {
	RenderTemplate(wr io.Writer, name string, data any) error
}

type pageRender struct {
	tmpls map[string]*template.Template
}

func NewPageRender(mapping map[string][]string) Render {
	var tmpls = make(map[string]*template.Template)

	for key, value := range mapping {
		t := template.Must(template.ParseFiles(value...))
		tmpls[key] = t
	}

	return &pageRender{tmpls}
}

func (r *pageRender) RenderTemplate(wr io.Writer, name string, data any) error {
	t, ok := r.tmpls[name]

	if !ok {
		return fmt.Errorf("missing template '%s'", name)
	}

	return t.ExecuteTemplate(wr, name, data)
}
