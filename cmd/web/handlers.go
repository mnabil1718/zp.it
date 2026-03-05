package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	urlib "net/url"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/mnabil1718/zp.it/internal/model"
	qrlib "github.com/mnabil1718/zp.it/internal/qr"
	"github.com/mnabil1718/zp.it/internal/shortener"
)

func (a *App) Health(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

type IndexData struct {
	Host string
}

func (a *App) Index(c *echo.Context) error {
	data := IndexData{Host: a.Config.Host}
	return c.Render(200, "index", data)
}

type Result struct {
	Short    string
	Original string
	QRCode   string // base64
}

func (a *App) Generate(c *echo.Context) error {
	url := c.FormValue("url")
	u, err := urlib.Parse(url)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid URL format")
	}

	// ensure host has a valid TLD
	host := u.Hostname()
	parts := strings.Split(host, ".")
	if len(parts) < 2 || parts[len(parts)-1] == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid URL format")
	}

	alias := c.FormValue("alias")
	qr := c.FormValue("qr") == "on"
	var code string

	if alias == "" {
		sc, err := shortener.Shorten(6)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to shorten URL")
		}

		code = sc
	} else {
		code = alias
	}

	if err := a.Models.Lookup.Insert(url, code); err != nil {
		if errors.Is(err, model.ErrAlreadyExists) {
			return echo.NewHTTPError(http.StatusConflict, "Code alias already exists")
		}

		fmt.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save lookup data")
	}

	s := a.Config.Host + code

	data := Result{
		Short:    s,
		Original: url,
	}

	if qr {
		png, err := qrlib.GenerateQR(s)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Cannot process QR generation")
		}

		data.QRCode = base64.StdEncoding.EncodeToString(png)
	}

	return c.Render(200, "result", data)
}

func (a *App) CodeHandler(c *echo.Context) error {
	cd := c.Param("code")

	lkp, err := a.Models.Lookup.GetByCode(cd)
	if err != nil {

		if errors.Is(err, model.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Short link is not found")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "Cannot lookup URL data")
	}

	return c.Redirect(http.StatusFound, lkp.Origin)
}
