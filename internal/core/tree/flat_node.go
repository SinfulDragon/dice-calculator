package tree

import "github.com/SinfulDragon/dice-calculator/internal/core/common"

type FlatNode struct {
	Raw   string
	Value int
	dice  []*common.Die
}

func (n *FlatNode) Evaluate() int {
	return n.Value
}

func (n *FlatNode) Roll() []*common.Die {
	if n.dice == nil {
		n.dice = []*common.Die{{Sides: 0, Value: n.Value}}
	}
	return n.dice
}

func (n *FlatNode) GetDice() []*common.Die {
	if n.dice == nil {
		n.dice = []*common.Die{{Sides: 0, Value: n.Value}}
	}
	return n.dice
}

func (n *FlatNode) Clone() FormulaNode {
	return &FlatNode{
		Raw:   n.Raw,
		Value: n.Value,
		dice:  n.dice,
	}
}

func (n *FlatNode) String() string {
	return n.Raw
}
