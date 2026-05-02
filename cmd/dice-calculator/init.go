package main

import (
	"github.com/SinfulDragon/dice-calculator/internal/core/modifiers/factory"
	"github.com/SinfulDragon/dice-calculator/internal/core/modifiers/reroll"
)

func Init() {
	// register modifiers
	factory.GlobalRegistry.Register("reroll", &reroll.RerollBuilder{})
}
