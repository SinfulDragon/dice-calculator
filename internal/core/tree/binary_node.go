package tree

import (
	"fmt"
	"slices"

	"github.com/SinfulDragon/dice-calculator/internal/core/common"
)

type BinaryOp int

const (
	BinaryPlus BinaryOp = iota
	BinaryMinus
	BinaryMul
	BinaryDiv
)

type BinaryNode struct {
	Raw   string
	Op    BinaryOp
	Left  FormulaNode
	Right FormulaNode
}

func (n *BinaryNode) Evaluate() int {
	switch n.Op {
	case BinaryPlus:
		return n.Left.Evaluate() + n.Right.Evaluate()
	case BinaryMinus:
		return n.Left.Evaluate() - n.Right.Evaluate()
	case BinaryMul:
		return n.Left.Evaluate() * n.Right.Evaluate()
	case BinaryDiv:
		return n.Left.Evaluate() / n.Right.Evaluate()
	default:
		return 0
	}
}

func (n *BinaryNode) Roll() []*common.Die {
	leftRolls := n.Left.Roll()
	rightRolls := n.Right.Roll()
	dice := slices.Concat(leftRolls, rightRolls)

	return dice
}

func (n *BinaryNode) GetDice() []*common.Die {
	leftDice := n.Left.GetDice()
	rightDice := n.Right.GetDice()
	dice := slices.Concat(leftDice, rightDice)

	return dice
}

func (n *BinaryNode) Clone() FormulaNode {
	return &BinaryNode{
		Raw:   n.Raw,
		Op:    n.Op,
		Left:  n.Left.Clone(),
		Right: n.Right.Clone(),
	}
}

func (n *BinaryNode) String() string {
	op := "+"
	switch n.Op {
	case BinaryMinus:
		op = "-"
	case BinaryMul:
		op = "*"
	case BinaryDiv:
		op = "/"
	}
	return fmt.Sprintf("(%s %s %s)", stringNode(n.Left), op, stringNode(n.Right))
}
