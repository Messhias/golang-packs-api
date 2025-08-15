package packages

import (
	"encoding/json"
	"os"
	"sort"

	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
	controller := NewPackageController()

	app.Get("/", controller.GetUI)

	app.Get("/pack-sizes", controller.PackSizes)

	// Update pack sizes
	app.Post("/pack-sizes", controller.PackSizes)

	// GET /calculate?items=251
	app.Get("/calculate", controller.GetCalculation)

	// POST /calculate { "items": 251 }
	app.Post("/calculate", controller.Calculate)

	app.Static("/", "./public")
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
