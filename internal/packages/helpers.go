package packages

import (
	"sort"
)

const MaxIntValue = 1 << 30

// greatestCommonDivisor helpers
func greatestCommonDivisor(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
func greatestCommonDivisorSlice(value []int) int {
	if len(value) == 0 {
		return 1
	}
	greatest := value[0]
	for i := 1; i < len(value); i++ {
		greatest = greatestCommonDivisor(greatest, value[i])
	}
	return greatest
}

func mountResult(bestTotal int, sizes []int, bestCounts []int) (Result, error) {
	res := Result{TotalItems: bestTotal, PacksUsed: map[int]int{}}
	for i, s := range sizes {
		if bestCounts[i] > 0 {
			res.PacksUsed[s] = bestCounts[i]
		}
	}
	return res, nil
}

func targetsAndCountsCalculator(
	target int,
	limit int,
	greatest int,
	try func(totalItems int) (int, []int, bool), bestTotal int, bestCounts []int) (int, []int, Result, error) {
	for tot := target; tot <= limit; tot += greatest {
		if _, counts, ok := try(tot); ok {
			bestTotal = tot
			bestCounts = counts
			break // the first one it's the fit  it is already the minor
		}
	}
	if bestTotal == -1 {
		return 0, nil, Result{}, NoCombinationError
	}

	return bestTotal, bestCounts, Result{}, nil
}

func reconstruction(totalItems int, sizes []int,
	previousStep []PreviousStepInfo, minimumPacksRequired []int) (int, []int, bool) {
	packCounts := make([]int, len(sizes))
	currentRemaining := totalItems
	for currentRemaining > 0 {
		currentPackSize := previousStep[currentRemaining].packSize
		if currentPackSize <= 0 {
			return 0, nil, false
		}

		getCurrentPackageIndex(sizes, currentPackSize, packCounts)
		currentRemaining = previousStep[currentRemaining].previousTotal
	}
	return minimumPacksRequired[totalItems], packCounts, true
}

// get the current package index  to be used in reconstruction
func getCurrentPackageIndex(sizes []int, currentPackSize int, packCounts []int) {
	for i, s := range sizes {
		if s == currentPackSize {
			packCounts[i]++
			break
		}
	}
}

func minimumPackageIterator(totalItems int, sizes []int) ([]int, []PreviousStepInfo, int, []int, bool, bool) {
	minimumPacksRequired := make([]int, totalItems+1)
	previousStep := make([]PreviousStepInfo, totalItems+1)
	for i := range minimumPacksRequired {
		minimumPacksRequired[i] = MaxIntValue
		previousStep[i] = PreviousStepInfo{-1, -1}
	}
	minimumPacksRequired[0] = 0

	for _, currentPackSize := range sizes {
		packageSizesControl(totalItems, currentPackSize, minimumPacksRequired, previousStep)
	}

	if minimumPacksRequired[totalItems] >= MaxIntValue {
		return nil, nil, 0, nil, false, true
	}
	return minimumPacksRequired, previousStep, 0, nil, false, false
}

func packageSizesControl(totalItems int, currentPackSize int, minimumPacksRequired []int, previousStep []PreviousStepInfo) {
	for currentTotal := currentPackSize; currentTotal <= totalItems; currentTotal++ {
		if minimumPacksRequired[currentTotal-currentPackSize]+1 < minimumPacksRequired[currentTotal] {
			minimumPacksRequired[currentTotal] = minimumPacksRequired[currentTotal-currentPackSize] + 1
			previousStep[currentTotal] = PreviousStepInfo{currentTotal - currentPackSize, currentPackSize}
		}
	}
}

// step 1: minor total of items >= orders.
// respecting the  greatestCommonDivisor
func mountTargets(items int, sizes []int) (int, int, Result, error, bool) {

	if items <= 0 {
		return 0, 0, Result{}, ErrInvalid, true
	}

	if len(sizes) == 0 {
		return 0, 0, Result{}, NoPackSizes, true
	}

	// order the sizes ( greater >= minor to pretty output).
	sort.Slice(sizes, func(i, j int) bool { return sizes[i] > sizes[j] })

	greatest := greatestCommonDivisorSlice(sizes)

	ceil := func(x, m int) int {
		if x%m == 0 {
			return x
		}
		return (x/m + 1) * m
	}

	return greatest, ceil(items, greatest), Result{}, nil, false
}
