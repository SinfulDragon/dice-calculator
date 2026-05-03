package models

import (
	"time"

	"github.com/SinfulDragon/dice-calculator/internal/core/common"
	"github.com/SinfulDragon/dice-calculator/internal/core/tree"
)

type RollResult struct {
	FormulaStr string
	Node       tree.FormulaNode
	Dice       []*common.Die
	Total      int
	Time       time.Time
}
