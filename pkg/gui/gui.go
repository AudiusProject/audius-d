package gui

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/AudiusProject/audius-d/pkg/logger"
	"github.com/labstack/echo/v4"
)

//go:embed dist/*
var embeddedFiles embed.FS

func StartGuiServer() error {
	e := echo.New()

	// Isolate 'ui/dist' as a separate segment in the embedded file system for efficient access
	subFS, err := fs.Sub(embeddedFiles, "dist")
	if err != nil {
		errmsg := fmt.Sprintf("failed to create sub filesystem: %v", err)
		return logger.Error(errmsg)
	}

	// Create a file server for the embedded files
	fileServer := http.FileServer(http.FS(subFS))

	// Serve static files directly
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", fileServer)))

	// Serve the index.html for all other GET requests
	e.GET("/*", func(c echo.Context) error {
		// Check if the requested file exists in the sub filesystem
		requestPath := c.Request().URL.Path
		if _, err := fs.Stat(subFS, strings.TrimPrefix(requestPath, "/")); err == nil {
			// If the file exists, let the file server handle the request
			fileServer.ServeHTTP(c.Response().Writer, c.Request())
			return nil
		}

		// Otherwise, serve index.html
		content, err := fs.ReadFile(subFS, "index.html")
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return c.HTMLBlob(http.StatusOK, content)
	})

	// Start the server with TLS (HTTPS)
	err = e.StartTLS(":2024", "deployments/tls/devnet-cert.pem", "deployments/tls/devnet-key.pem")
	if err != nil {
		return logger.Error(err)
	}
	return nil
}

func main() {
	StartGuiServer()
}
