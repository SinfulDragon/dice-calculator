package tree

import (
	"fmt"

	"github.com/SinfulDragon/dice-calculator/internal/core/common"
)

type DiceNode struct {
	Raw   string
	Count int
	Sides int
	dice  []*common.Die
}

func (n *DiceNode) Evaluate() int {
	if n.dice == nil {
		return 0
	}
	sum := 0
	for _, die := range n.dice {
		sum += die.Value
	}
	return sum
}

func (n *DiceNode) Roll() []*common.Die {

	if n.dice == nil || cap(n.dice) < n.Count {
		n.dice = make([]*common.Die, n.Count)
	} else {
		n.dice = n.dice[:n.Count]
	}

	for i := range n.dice {
		if n.dice[i] == nil {
			n.dice[i] = common.NewDie(n.Sides)
		} else {
			n.dice[i].Roll()
		}
	}

	return n.dice
}

func (n *DiceNode) GetDice() []*common.Die {
	return n.dice
}

func (n *DiceNode) Clone() FormulaNode {
	return &DiceNode{
		Raw:   n.Raw,
		Count: n.Count,
		Sides: n.Sides,
	}
}

func (n *DiceNode) String() string {
	if n.Count == 1 {
		return fmt.Sprintf("d%d", n.Sides)
	}
	return fmt.Sprintf("%dd%d", n.Count, n.Sides)
}
