package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/mnabil1718/zp.it/internal/shortener"
)

type Result struct {
	Short    string
	Original string
}

func Health(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

func Index(c *echo.Context) error {
	return c.Render(200, "index", nil)
}

func Generate(c *echo.Context) error {

	fmt.Println(c.FormValue("url"))

	url := c.FormValue("url")
	qr := c.FormValue("qr") == "on"

	code, err := shortener.Shorten(6)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Something went wrong")
	}

	// TODO: wire qr later
	_ = qr

	data := Result{
		Short:    "http://localhost:8080/" + code,
		Original: url,
	}

	return c.Render(200, "result", data)
}
