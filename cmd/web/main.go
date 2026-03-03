package main

import (
	"html/template"
	"io"
	"os"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/mnabil1718/zp.it/internal/db"
	"github.com/mnabil1718/zp.it/internal/model"
)

type App struct {
	Models *model.Models
}

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
	dbp := os.Getenv("DB_PATH")
	if dbp == "" {
		dbp = "data.db"
	}

	reset := os.Getenv("APP_ENV") != "production"

	db := db.NewSQLiteDB(dbp, reset)
	defer db.Close()

	lu := model.NewSQliteLookup(db)
	models := model.NewModels(lu)

	app := App{
		Models: models,
	}

	e := echo.New()
	e.Renderer = newTemplate()
	e.Static("/static", "static")
	e.Use(middleware.RequestLogger())
	e.HTTPErrorHandler = ErrorHandler

	e.GET("/", app.Index)
	e.GET("/health", app.Health)
	e.GET("/:code", app.CodeHandler)
	e.POST("/generate", app.Generate)

	if err := e.Start(":8080"); err != nil {
		e.Logger.Error("Failed to start server", "error", err)
	}
}
