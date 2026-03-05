package main

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/mnabil1718/zp.it/internal/cache"
	"github.com/mnabil1718/zp.it/internal/config"
	"github.com/mnabil1718/zp.it/internal/db"
	"github.com/mnabil1718/zp.it/internal/model"
)

type App struct {
	Models *model.Models
	Config *config.Config
	Server *http.Server
}

func NewApp(cfg *config.Config, db *sql.DB, cache cache.ICache) *App {
	lu := model.NewSQliteLookup(db, cache)
	models := model.NewModels(lu)
	return &App{
		Models: models,
		Config: cfg,
	}
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

func (app *App) setupServer() {
	e := echo.New()
	e.Renderer = newTemplate()
	e.Static("/static", "static")
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	e.Use(middleware.RequestLogger())
	e.HTTPErrorHandler = ErrorHandler

	e.GET("/", app.Index)
	e.GET("/health", app.Health)
	e.GET("/:code", app.CodeHandler)
	e.POST("/generate", app.Generate)

	app.Server = &http.Server{
		Addr:    fmt.Sprintf(":%d", app.Config.Port),
		Handler: e,
	}
}

func (app *App) serve() {

	go func() {
		slog.Info("server is starting on", "port", app.Config.Port, "env", app.Config.Env)
		if err := app.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.Server.Shutdown(ctx); err != nil {
		slog.Error("shutting down server in 10s", "error", err)
	}
}

func bootstrap() {
	cfg := config.Load()

	dbp := cfg.DBPath
	if dbp == "" {
		dbp = "data.db"
	}
	reset := cfg.Env != "production"

	db := db.NewSQLiteDB(dbp, reset)
	defer db.Close()

	rdb := cache.NewRedisClient(cfg)
	defer rdb.Close()

	app := NewApp(cfg, db, rdb)
	app.setupServer()
	app.serve()
}
