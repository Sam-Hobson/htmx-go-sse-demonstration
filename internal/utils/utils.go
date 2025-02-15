package utils

import (
	"bytes"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

var Template = &templateRenderer{
	templates: template.Must(template.ParseGlob("public/views/*.html")),
}

type templateRenderer struct {
	templates *template.Template
}

func (t *templateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func RenderTemplateToString(name string, data interface{}) (string, error) {
	var buf bytes.Buffer
	if err := Template.templates.ExecuteTemplate(&buf, name, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
