package tree

import (
	"github.com/SinfulDragon/dice-calculator/internal/core/common"
)

type ModifierNode struct {
	Raw      string
	Modifier common.Modifier
	Child    FormulaNode
}

func (n *ModifierNode) Evaluate() int {
	return n.Child.Evaluate()
}

func (n *ModifierNode) Roll() []*common.Die {
	rolls := n.Child.Roll()
	n.Modifier.Apply(rolls)
	return rolls
}

func (n *ModifierNode) GetDice() []*common.Die {
	return n.Child.GetDice()
}

func (n *ModifierNode) Clone() FormulaNode {
	return &ModifierNode{
		Raw:      n.Raw,
		Modifier: n.Modifier,
		Child:    n.Child.Clone(),
	}
}

func (n *ModifierNode) String() string {
	return n.Raw
}
