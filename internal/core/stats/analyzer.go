package stats

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/SinfulDragon/dice-calculator/internal/core/tree"
)

type Analyzer struct {
	formula tree.FormulaNode
}

func NewAnalyzer(formula tree.FormulaNode) *Analyzer {
	return &Analyzer{formula: formula}
}

func (a *Analyzer) MonteCarlo(iterations int) (Distribution, error) {
	if iterations <= 0 {
		return Distribution{}, fmt.Errorf("iterations must be positive, got %d", iterations)
	}
	if a.formula == nil {
		return Distribution{}, fmt.Errorf("formula is nil")
	}

	workers := min(runtime.GOMAXPROCS(0), iterations)

	perWorker := iterations / workers
	remainder := iterations % workers

	type result struct {
		local map[int]int
	}

	results := make(chan result, workers)
	var wg sync.WaitGroup

	for w := range workers {
		count := perWorker
		if w < remainder {
			count++
		}
		if count == 0 {
			continue
		}

		wg.Add(1)
		go func(batchSize int) {
			defer wg.Done()

			f := a.formula.Clone()
			local := make(map[int]int, 64)

			for range batchSize {
				f.Roll()
				v := f.Evaluate()
				local[v]++
			}

			results <- result{local: local}
		}(count)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	outcomes := make(map[int]int, 128)
	total := 0

	for r := range results {
		for v, c := range r.local {
			outcomes[v] += c
			total += c
		}
	}

	return Distribution{
		outcomes: outcomes,
		total:    total,
	}, nil
}
