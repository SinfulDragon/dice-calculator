package stats

import (
	"math"
	"testing"

	"github.com/SinfulDragon/dice-calculator/internal/core/tree"
)

func TestAnalyzer_MonteCarlo_InvalidInput(t *testing.T) {
	_, err := NewAnalyzer(nil).MonteCarlo(100)
	if err == nil {
		t.Fatal("expected error for nil formula")
	}

	node := &tree.DiceNode{Raw: "1d6", Count: 1, Sides: 6}
	_, err = NewAnalyzer(node).MonteCarlo(0)
	if err == nil {
		t.Fatal("expected error for zero iterations")
	}

	_, err = NewAnalyzer(node).MonteCarlo(-1)
	if err == nil {
		t.Fatal("expected error for negative iterations")
	}
}

func TestAnalyzer_MonteCarlo_2d6(t *testing.T) {
	node := &tree.DiceNode{Raw: "2d6", Count: 2, Sides: 6}
	analyzer := NewAnalyzer(node)

	dist, err := analyzer.MonteCarlo(200_000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	summary := dist.Summary()

	if summary.Min != 2 {
		t.Errorf("min = %v, want 2", summary.Min)
	}
	if summary.Max != 12 {
		t.Errorf("max = %v, want 12", summary.Max)
	}

	expectedMean := 7.0
	if math.Abs(summary.Mean-expectedMean) > 0.1 {
		t.Errorf("mean = %v, want ~%v", summary.Mean, expectedMean)
	}

	// Check that all outcomes 2..12 have non-zero count
	for i := 2; i <= 12; i++ {
		if dist.Count(i) == 0 {
			t.Errorf("outcome %d has zero count", i)
		}
	}

	// Probability of 7 should be highest
	maxProb := 0.0
	maxVal := 0
	for _, v := range dist.Outcomes() {
		p := dist.Probability(v)
		if p > maxProb {
			maxProb = p
			maxVal = v
		}
	}
	if maxVal != 7 {
		t.Errorf("highest probability at %d, want 7", maxVal)
	}
}

func TestDistribution_Summary_Empty(t *testing.T) {
	d := Distribution{}
	s := d.Summary()
	if s != (Summary{}) {
		t.Errorf("expected zero Summary, got %+v", s)
	}
}
