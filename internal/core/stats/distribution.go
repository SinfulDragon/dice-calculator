package stats

import (
	"math"
	"sort"
)

type Distribution struct {
	outcomes map[int]int
	total    int
}

func (d Distribution) Count(value int) int {
	if d.outcomes == nil {
		return 0
	}
	return d.outcomes[value]
}

func (d Distribution) Probability(value int) float64 {
	if d.total == 0 {
		return 0
	}
	return float64(d.Count(value)) / float64(d.total)
}

func (d Distribution) Outcomes() []int {
	if d.outcomes == nil {
		return nil
	}
	keys := make([]int, 0, len(d.outcomes))
	for k := range d.outcomes {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

func (d Distribution) Total() int {
	return d.total
}

type Summary struct {
	Min    float64
	Max    float64
	Mean   float64
	StdDev float64
}

func (d Distribution) Summary() Summary {
	if d.total == 0 {
		return Summary{}
	}

	var min, max int
	var sum, sumSq float64
	first := true

	for value, count := range d.outcomes {
		if first || value < min {
			min = value
		}
		if first || value > max {
			max = value
		}
		first = false

		fv := float64(value)
		fc := float64(count)
		sum += fv * fc
		sumSq += fv * fv * fc
	}

	mean := sum / float64(d.total)
	variance := sumSq/float64(d.total) - mean*mean
	if variance < 0 {
		variance = 0
	}

	return Summary{
		Min:    float64(min),
		Max:    float64(max),
		Mean:   mean,
		StdDev: math.Sqrt(variance),
	}
}
