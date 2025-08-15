package packages

import (
	"log"
	"net/http"
	"sort"

	"github.com/gofiber/fiber/v2"
)

type PackageController struct {
	packs []int
}

func NewPackageController() *PackageController {
	packs := loadPacksConfigurations()
	return &PackageController{
		packs: packs,
	}
}

func loadPacksConfigurations() []int {
	packs, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}
	return packs
}

func (pc *PackageController) GetUI(c *fiber.Ctx) error {
	return c.SendFile("./public/index.html")
}

func (pc *PackageController) PackSizes(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"pack_sizes": pc.packs,
	})
}

func (pc *PackageController) UpdatePackSizes(c *fiber.Ctx) error {
	var body struct {
		PackSizes []int `json:"pack_sizes"`
	}

	packs := pc.packs

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
}

func (pc *PackageController) GetCalculation(c *fiber.Ctx) error {
	items, err := c.QueryInt("items", 0), error(nil)
	if items <= 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "items must be > 0"})
	}
	res, err := RetrievePackages(items, pc.packs)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (pc *PackageController) Calculate(c *fiber.Ctx) error {
	var body struct {
		Items int `json:"items"`
	}
	if err := c.BodyParser(&body); err != nil || body.Items <= 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	res, err := RetrievePackages(body.Items, pc.packs)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}
