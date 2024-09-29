package handler

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

type Template struct {
	*template.Template
}

func (t Template) Render(w io.Writer, name string, data interface{}, _ echo.Context) error {
	return t.ExecuteTemplate(w, name, data)
}
