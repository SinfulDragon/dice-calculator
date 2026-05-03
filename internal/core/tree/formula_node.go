package tree

import "github.com/SinfulDragon/dice-calculator/internal/core/common"

type FormulaNode interface {
	Evaluate() int
	Roll() []*common.Die
	GetDice() []*common.Die
	Clone() FormulaNode
}

func stringNode(n FormulaNode) string {
	if s, ok := n.(interface{ String() string }); ok {
		return s.String()
	}
	return ""
}
