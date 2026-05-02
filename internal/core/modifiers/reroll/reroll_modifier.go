package reroll

import (
	"cmp"
	"slices"

	"github.com/SinfulDragon/dice-calculator/internal/core/common"
)

type RerollMode int

const (
	JustReroll RerollMode = iota
	RerollHighest
	RerollLowest
	RerollBelow
	RerollAbove
	RerollExact
)

type RerollModifier struct {
	Mode     RerollMode
	Values   []int // values to reroll when in RerollExact mode
	Limit    int   // number of dice to reroll (0 = all)
	MinValue int   // threshold value for rerolling (used in RerollBelow and RerollAbove modes)
}

func (r *RerollModifier) Apply(dice []*common.Die) {
	switch r.Mode {
	case JustReroll:
		r.justReroll(dice)
	case RerollHighest:
		r.rerollHighest(dice)
	case RerollLowest:
		r.rerollLowest(dice)
	case RerollBelow:
		r.rerollBelow(dice)
	case RerollAbove:
		r.rerollAbove(dice)
	case RerollExact:
		r.rerollExact(dice)
	default:
		r.justReroll(dice)
	}
}

func (r *RerollModifier) checkLimit(i int) bool {
	if r.Limit == 0 {
		return true
	}
	return i < r.Limit
}

func (r *RerollModifier) justReroll(dice []*common.Die) {
	for i, die := range dice {
		if r.checkLimit(i) {
			die.Roll()
		}
	}
}

func (r *RerollModifier) rerollHighest(dice []*common.Die) {
	sorted := slices.SortedStableFunc(slices.Values(dice), func(a, b *common.Die) int {
		return cmp.Compare(b.Value, a.Value)
	})
	for i, die := range sorted {
		if r.checkLimit(i) {
			die.Roll()
		}
	}
}

func (r *RerollModifier) rerollLowest(dice []*common.Die) {
	sorted := slices.SortedStableFunc(slices.Values(dice), func(a, b *common.Die) int {
		return cmp.Compare(a.Value, b.Value)
	})
	for i, die := range sorted {
		if r.checkLimit(i) {
			die.Roll()
		}
	}
}

func (r *RerollModifier) rerollBelow(dice []*common.Die) {
	sorted := slices.SortedStableFunc(slices.Values(dice), func(a, b *common.Die) int {
		return cmp.Compare(a.Value, b.Value)
	})
	count := 0
	for _, die := range sorted {
		if die.Value < r.MinValue && r.checkLimit(count) {
			die.Roll()
			count++
		}
	}
}

func (r *RerollModifier) rerollAbove(dice []*common.Die) {
	sorted := slices.SortedStableFunc(slices.Values(dice), func(a, b *common.Die) int {
		return cmp.Compare(b.Value, a.Value)
	})
	count := 0
	for _, die := range sorted {
		if die.Value > r.MinValue && r.checkLimit(count) {
			die.Roll()
			count++
		}
	}
}

func (r *RerollModifier) rerollExact(dice []*common.Die) {
	count := 0
	for _, die := range dice {
		if slices.Contains(r.Values, die.Value) && r.checkLimit(count) {
			die.Roll()
			count++
		}
	}
}
