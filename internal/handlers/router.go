package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRouter(app *fiber.App) {
	app.Get("/health", HealthHandler)
	app.Get("/", HomeHandler)
	app.Get("/artist", ArtistHandler)

	app.Use(func(c *fiber.Ctx) error {
		return showErrorPage(c, fiber.StatusNotFound, "404: Page Not Found")
	})
}
