package main

import (
	"github.com/labstack/echo/v4"
)

func main() {
	// Initialize Echo
	e := echo.New()

	e.POST("/auth", nil)
	// Start server
	e.Logger.Fatal(e.Start(":8000"))
}
