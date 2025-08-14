package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"pack-calculator/internal/solver"
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
		return c.SendString(indexHTML)
	})

	// GET /calculate?items=251
	app.Get("/calculate", func(c *fiber.Ctx) error {
		items, err := c.QueryInt("items", 0), error(nil)
		if items <= 0 {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "items must be > 0"})
		}
		res, err := solver.Solve(items, packs)
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
		res, err := solver.Solve(body.Items, packs)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(res)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}

const indexHTML = `<!doctype html>
<html lang="en"><head>
<meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>Pack Calculator</title>
<style>
body{font-family:system-ui,Arial,sans-serif;margin:40px}
input,button{font-size:16px;padding:8px 10px}
.card{border:1px solid #ddd;border-radius:12px;padding:20px;max-width:520px}
.result{margin-top:16px;white-space:pre;font-family:ui-monospace,Consolas,monospace}
</style>
</head><body>
<h1>Pack Calculator</h1>
<div class="card">
  <label>Items: <input id="items" type="number" min="1" value="251"></label>
  <button id="calc">Calculate</button>
  <div class="result" id="out"></div>
</div>
<script>
document.getElementById('calc').onclick = async () => {
  const items = document.getElementById('items').value;
  const res = await fetch('/calculate?items='+items);
  const data = await res.json();
  if(!res.ok){ document.getElementById('out').textContent = 'Error: '+(data.error||res.status); return; }
  let lines = [];
  lines.push('Total items sent: '+data.total_items);
  const packs = Object.entries(data.packs_used).sort((a,b)=>Number(b[0])-Number(a[0]));
  for(const [size,count] of packs){
  }
  document.getElementById('out').textContent = lines.join('\n');
};
</script>
</body></html>`
