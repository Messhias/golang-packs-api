package packages

type PreviousStepInfo struct{ previousTotal, packSize int }

type Result struct {
	TotalItems int         `json:"total_items"`
	PacksUsed  map[int]int `json:"packs_used"` // size -> count
}
