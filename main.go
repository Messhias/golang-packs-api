package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"pack-calculator/internal/packages"
	"sort"

	"github.com/gofiber/fiber/v2"
)

type Config struct {
	Packs []int `json:"packs"`
}

func loadConfig() ([]int, error) {
	f, err := os.Open("packs.json")
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}

	// ensure the order
	sort.Slice(cfg.Packs, func(i, j int) bool { return cfg.Packs[i] > cfg.Packs[j] })
	return cfg.Packs, nil
}

func main() {
	packs, err := loadConfig()
	if err != nil {
		log.Fatalf("failed to load packs.json: %v", err)
	}

	app := fiber.New()

	// UI
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./public/index.html")
	})

	// Get current pack sizes
	app.Get("/pack-sizes", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"pack_sizes": packs,
		})
	})

	// Update pack sizes
	app.Post("/pack-sizes", func(c *fiber.Ctx) error {
		var body struct {
			PackSizes []int `json:"pack_sizes"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		// Validate pack sizes
		for _, size := range body.PackSizes {
			if size <= 0 {
				return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Pack sizes must be positive integers"})
			}
		}

		// Update the packs slice
		packs = body.PackSizes
		sort.Slice(packs, func(i, j int) bool { return packs[i] > packs[j] })

		return c.JSON(fiber.Map{
			"pack_sizes": packs,
		})
	})

	// GET /calculate?items=251
	app.Get("/calculate", func(c *fiber.Ctx) error {
		items, err := c.QueryInt("items", 0), error(nil)
		if items <= 0 {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "items must be > 0"})
		}
		res, err := packages.RetrievePackages(items, packs)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(res)
	})

	// POST /calculate { "items": 251 }
	app.Post("/calculate", func(c *fiber.Ctx) error {
		var body struct {
			Items int `json:"items"`
		}
		if err := c.BodyParser(&body); err != nil || body.Items <= 0 {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
		}
		res, err := packages.RetrievePackages(body.Items, packs)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(res)
	})

	app.Static("/", "./public")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}
