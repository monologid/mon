package mon

import (
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type IEcho interface {
	// GetApp returns echo.Echo object
	GetApp() *echo.Echo

	// Start starts the application server
	Start(addr string)
}

type Echo struct {
	App *echo.Echo
}

func (e *Echo) GetApp() *echo.Echo {
	return e.App
}

func (e *Echo) Start(addr string) {
	if len(addr) == 0 {
		e.App.Logger.Fatal(e.App.Start(":1323"))
	}

	e.App.Logger.Fatal(e.App.Start(addr))
}

func NewEcho() IEcho {
	e := &Echo{App: echo.New()}

	e.App.Use(middleware.Logger())
	e.App.Use(middleware.Recover())
	e.App.Use(middleware.CSRF())
	e.App.Use(middleware.Secure())

	prom := prometheus.NewPrometheus("echo", nil)
	prom.Use(e.App)

	return e
}
