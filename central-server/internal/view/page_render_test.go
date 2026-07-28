package view

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPageRender_RenderTemplate_Success(t *testing.T) {

	dir := t.TempDir()
	templatePath := filepath.Join(dir, "home.html")
	err := os.WriteFile(
		templatePath,
		[]byte(`
			<html>
				<body>
					<h1>{{.Title}}</h1>
				</body>
			</html>
		`),
		0644,
	)
	assert.NoError(t, err)

	renderer := NewPageRender(map[string][]string{
		"home.html": {
			templatePath,
		},
	})

	var buffer bytes.Buffer
	err = renderer.RenderTemplate(&buffer, "home.html", map[string]string{"Title": "Hello World"})

	assert.NoError(t, err)

	assert.Contains(t, buffer.String(), "Hello World")
}

func TestPageRender_RenderTemplate_MissingTemplate(t *testing.T) {

	renderer := NewPageRender(map[string][]string{})
	var buffer bytes.Buffer
	err := renderer.RenderTemplate(&buffer, "missing.html", nil)

	assert.Error(t, err)
	assert.EqualError(t, err, "missing template 'missing.html'")
}

func TestPageRender_RenderTemplate_InvalidData(t *testing.T) {

	dir := t.TempDir()
	templatePath := filepath.Join(dir, "page.html")

	err := os.WriteFile(
		templatePath,
		[]byte(`
			{{.User.Name}}
		`),
		0644,
	)

	assert.NoError(t, err)

	renderer := NewPageRender(map[string][]string{"page.html": {templatePath}})
	var buffer bytes.Buffer
	err = renderer.RenderTemplate(&buffer, "page.html", "wrong data")

	assert.Error(t, err)
}

func TestPageRender_MultipleTemplates(t *testing.T) {

	dir := t.TempDir()
	header := filepath.Join(dir, "header.html")
	page := filepath.Join(dir, "page.html")

	assert.NoError(t, os.WriteFile(header, []byte(`HEADER`), 0644))
	assert.NoError(t, os.WriteFile(page, []byte(`PAGE`), 0644))

	renderer := NewPageRender(map[string][]string{"page.html": {header, page}})
	var buffer bytes.Buffer
	err := renderer.RenderTemplate(&buffer, "page.html", nil)

	assert.NoError(t, err)
	assert.Contains(t, buffer.String(), "PAGE")
}
