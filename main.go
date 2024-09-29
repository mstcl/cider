package main

import (
	"embed"
	"flag"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mstcl/cider/internal/handler"
)

var addr = flag.String("addr", "0.0.0.0:8080", "address to listen on")

//go:embed web
var web embed.FS

func main() {
	flag.Parse()

	r := echo.New()
	r.Use(middleware.Logger())
	r.Use(middleware.Recover())
	r.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		HTML5:      true,
		Root:       "web",
		Filesystem: http.FS(web),
	}))
	r.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	r.GET("/", handler.Index)

	t := handler.Template{Template: template.Must(template.ParseFS(web, "web/templates/index.tmpl"))}
	r.Renderer = t

	r.IPExtractor = echo.ExtractIPDirect()

	r.Server.ReadHeaderTimeout = 5 * time.Second
	r.Server.ReadTimeout = 5 * time.Second
	r.Server.WriteTimeout = 5 * time.Second
	r.Server.IdleTimeout = 60 * time.Second

	if err := r.Start(*addr); err != http.ErrServerClosed {
		log.Fatal("http error:", err)
	}
}
