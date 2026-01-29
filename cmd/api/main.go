package main

import (
	"groupie-tracker-visualizations/internal/handlers"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	engine := html.New("./web/templates", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/static", "./web/static")

	handlers.SetupRouter(app)

	log.Println("Server started at http://localhost:8000")
	if err := app.Listen(":8000"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
