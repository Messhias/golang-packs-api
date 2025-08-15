package main

import (
	"log"
	"os"
	"pack-calculator/internal/packages"

	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()

	packages.Routes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}
