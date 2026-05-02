package common

import (
	"fmt"
	"math/rand/v2"
)

type Die struct {
	Sides int
	Value int
}

func NewDie(sides int) *Die {
	die := &Die{Sides: sides}
	die.Roll()
	return die
}

func (d *Die) Roll() int {
	d.Value = rand.N(d.Sides) + 1
	return d.Value
}

func (d *Die) String() string {
	return fmt.Sprintf("d%d=%d", d.Sides, d.Value)
}
