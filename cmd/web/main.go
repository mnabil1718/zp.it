package main

import (
	"html/template"
	"io"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(c *echo.Context, w io.Writer, name string, data any) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Template {
	funcs := template.FuncMap{
		"stripScheme": func(url string) string {
			url = strings.TrimPrefix(url, "https://")
			url = strings.TrimPrefix(url, "http://")
			return url
		},
	}
	return &Template{
		templates: template.Must(template.New("").Funcs(funcs).ParseGlob("ui/*.html")),
	}
}

func main() {
	e := echo.New()
	e.Renderer = newTemplate()
	e.Static("/static", "static")
	e.Use(middleware.RequestLogger())

	e.GET("/health", Health)
	e.GET("/", Index)
	e.POST("/generate", Generate)

	if err := e.Start(":8080"); err != nil {
		e.Logger.Error("Failed to start server", "error", err)
	}
}
