package packages

func RetrievePackages(items int, sizes []int) (Result, error) {
	greatest, targetResult, result, err, done := mountTargets(items, sizes)
	if done {
		return result, err
	}
	target := targetResult

	// step 2: for the minor total as possible, minize the number of the
	// packages (unlimited changes); if it's not possible to mount the
	// 'target' (impossible minor values), we increase in greatest until find the matches.
	maxSize := sizes[0]
	limit := items + maxSize*10 // usually we found before.

	runtimeFunction := func(totalItems int) (int, []int, bool) {
		minimumPacksRequired, previousStep, i, integers, b, done2 := minimumPackageIterator(totalItems, sizes)
		if done2 {
			return i, integers, b
		}

		// reconstruct.
		return reconstruction(totalItems, sizes, previousStep, minimumPacksRequired)
	}
	try := runtimeFunction

	bestTotal := -1
	var bestCounts []int

	bestTotal, bestCounts, r, err2 := targetsAndCountsCalculator(target, limit, greatest, try, bestTotal, bestCounts)
	if err2 != nil {
		return r, err2
	}

	return mountResult(bestTotal, sizes, bestCounts)
}
