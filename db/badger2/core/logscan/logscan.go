package logscan

import "github.com/ethereum/go-ethereum/crypto"

// N = 1e9 (Expected number of entries)
// P = 0.001 (Desired false-positive probability)
const HashSize = 5 // ceil(ln(N / P) / ln(256))

func BitmaskContains(bitmask, v int) bool {
	return ((1 << v) & bitmask) > 0
}

func CalcHash(features [][]byte, bitmask int) []byte {
	hash := crypto.NewKeccakState()
	for i := 0; i < len(features); i++ {
		if BitmaskContains(bitmask, i) {
			hash.Write(features[i])
		}
	}
	hashBytes := make([]byte, HashSize)
	hash.Read(hashBytes)
	return hashBytes
}

func SelectSearchBitmask(featureFilters [][][]byte, maxIterators uint) int {
	maxBitmask := (1 << len(featureFilters)) - 1
	bestBitmask, bestFiltersCnt, bestIteratorsCnt := 0, 0, 1
	for bitmask := 1; bitmask <= maxBitmask; bitmask++ {
		valid, filtersCnt, iteratorsCnt := validateSearchBitmask(featureFilters, maxIterators, bitmask)
		if !valid {
			continue
		}
		if filtersCnt < bestFiltersCnt {
			continue
		}
		if filtersCnt == bestFiltersCnt && iteratorsCnt > bestIteratorsCnt {
			continue
		}
		bestBitmask, bestFiltersCnt, bestIteratorsCnt = bitmask, filtersCnt, iteratorsCnt
	}
	return bestBitmask
}

func validateSearchBitmask(featureFilters [][][]byte, maxIterators uint, bitmask int) (bool, int, int) {
	filtersCnt, iteratorsCnt := 0, 1
	for i, filter := range featureFilters {
		if !BitmaskContains(bitmask, i) {
			continue
		}
		if len(filter) == 0 || uint64(iteratorsCnt)*uint64(len(filter)) > uint64(maxIterators) {
			return false, 0, 0
		}
		filtersCnt++
		iteratorsCnt *= len(filter)
	}
	return true, filtersCnt, iteratorsCnt
}

func GenerateSearchHashes(featureFilters [][][]byte, bitmask int) map[string]struct{} {
	result := make(map[string]struct{})
	generateSearсhHashes(featureFilters, make([][]byte, len(featureFilters)), 0, bitmask, result)
	return result
}

func generateSearсhHashes(featureFilters [][][]byte, selectedFeatures [][]byte, curFilter, bitmask int, result map[string]struct{}) {
	if curFilter == len(featureFilters) {
		result[string(CalcHash(selectedFeatures, bitmask))] = struct{}{}
		return
	}
	if !BitmaskContains(bitmask, curFilter) {
		generateSearсhHashes(featureFilters, selectedFeatures, curFilter+1, bitmask, result)
		return
	}
	for _, selectedFeatures[curFilter] = range featureFilters[curFilter] {
		generateSearсhHashes(featureFilters, selectedFeatures, curFilter+1, bitmask, result)
	}
}
