package main

import (
	"audioconverter/src/converter"
	"audioconverter/utils/healthcheck"

	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)
const port = ":3000"

func main() {
    
	app := fiber.New()
	api := app.Group("/api")
	v1 := api.Group("/v1")
	v2 := api.Group("/v2")

	app.Get("/swagger/*", static.New("./public/swagger"))

	app.Get("/swagger.json", func(c fiber.Ctx) error {
		return c.SendFile("./docs/swagger.json")
	})
	
	app.Get("/health", healthcheck.HealthCheckHeandler)

	v1.Post("/convert", converter.ConverHeandler)
	v2.Post("/convert", converter.ConverV2Heandler)
	log.Fatal(app.Listen(port))
}