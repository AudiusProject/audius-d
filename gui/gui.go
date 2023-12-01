package gui

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

type Server struct {
	quit chan os.Signal
}

//go:embed ui/dist/*
var embeddedFiles embed.FS

func StartGuiServer() {
	e := echo.New()

	// Create a sub-filesystem for the 'ui/dist' directory
	subFS, err := fs.Sub(embeddedFiles, "ui/dist")
	if err != nil {
		log.Fatalf("failed to create sub filesystem: %v", err)
	}

	// Create a file server for the embedded files
	fileServer := http.FileServer(http.FS(subFS))

	// Route to serve index.html at /
	e.GET("/", func(c echo.Context) error {
		// Serve the embedded "index.html" file
		content, err := fs.ReadFile(subFS, "index.html")
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return c.HTMLBlob(http.StatusOK, content)
	})

	// Serve other static files
	e.GET("/*", echo.WrapHandler(http.StripPrefix("/", fileServer)))

	e.Logger.Fatal(e.Start(":2024"))
}
