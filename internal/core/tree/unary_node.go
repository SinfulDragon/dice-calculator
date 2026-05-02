package tree

import (
	"github.com/SinfulDragon/dice-calculator/internal/core/common"
)

type UnaryOp int

const (
	UnaryPlus UnaryOp = iota
	UnaryMinus
)

type UnaryNode struct {
	Raw   string
	Op    UnaryOp
	Child FormulaNode
}

func (n *UnaryNode) Evaluate() int {
	switch n.Op {
	case UnaryPlus:
		return n.Child.Evaluate()
	case UnaryMinus:
		return -n.Child.Evaluate()
	default:
		return 0
	}
}

func (n *UnaryNode) Roll() []*common.Die {
	return n.Child.Roll()
}

func (n *UnaryNode) GetDice() []*common.Die {
	return n.Child.GetDice()
}

func (n *UnaryNode) String() string {
	return n.Raw
}
