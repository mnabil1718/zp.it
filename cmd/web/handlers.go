package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	urlib "net/url"
	"strings"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/mnabil1718/zp.it/internal/model"
	qrlib "github.com/mnabil1718/zp.it/internal/qr"
	"github.com/mnabil1718/zp.it/internal/shortener"
)

func (a *App) Health(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

func (a *App) Index(c *echo.Context) error {
	return c.Render(http.StatusOK, "index", map[string]any{
		"Host": a.Config.Host,
	})
}

func (a *App) Counter(c *echo.Context) error {
	return c.Render(http.StatusOK, "counter", nil)
}

type Result struct {
	Short    string
	Original string
	QRCode   string // base64
}

func (a *App) Generate(c *echo.Context) error {
	url := c.FormValue("url")
	url = strings.Trim(url, " ")
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

type CounterResult struct {
	Clicks int
	Since  time.Time
}

func (a *App) GetCounterData(c *echo.Context) error {
	// NOTE: user might input only code or whole short link
	url := c.FormValue("url")
	url = strings.Trim(url, " ")

	var code string
	// Check if full URL or just code
	if strings.Contains(url, "://") || strings.Contains(url, ".") {
		u, err := urlib.Parse(url)
		if err != nil || (u.Host != "" && (u.Path == "" || u.Path == "/")) {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid URL format")
		}

		code = strings.TrimPrefix(u.Path, "/")
	} else {
		code = url
	}

	lkp, err := a.Models.Lookup.GetByCode(code)
	if err != nil {

		if errors.Is(err, model.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Short link is not found")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "Cannot lookup URL data")
	}

	data := CounterResult{
		Clicks: lkp.Clicks,
		Since:  lkp.CreatedAt,
	}

	c.Response().Header().Set("Cache-Control", "no-store")
	return c.Render(200, "counter-result", data)
}

func (a *App) CodeHandler(c *echo.Context) error {
	cd := c.Param("code")

	origin, err := a.Models.Lookup.GetOriginByCode(cd)
	if err != nil {

		if errors.Is(err, model.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Short link is not found")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "Cannot lookup URL data")
	}

	// runs in goroutine, non-blocking redirect
	go func() {
		if err := a.Models.Lookup.IncrementClicks(cd); err != nil {
			slog.Error("Cannot increment clicks", "code", cd, "error", err)
		}
	}()

	return c.Redirect(http.StatusFound, origin)
}
