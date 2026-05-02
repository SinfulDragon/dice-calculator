package tree

import "github.com/SinfulDragon/dice-calculator/internal/core/common"

type FormulaNode interface {
	Evaluate() int
	Roll() []*common.Die
	GetDice() []*common.Die
}
