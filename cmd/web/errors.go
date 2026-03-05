package main

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func ErrorHandler(c *echo.Context, err error) {

	code := http.StatusInternalServerError
	msg := err.Error()
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message
	}

	htmxReq := c.Request().Header.Get("HX-Request")
	if htmxReq == "true" {
		c.Render(code, "error-message", msg)
	} else {
		c.Render(code, "error-page", map[string]any{
			"Code":  code,
			"Error": msg,
		})
	}

}
