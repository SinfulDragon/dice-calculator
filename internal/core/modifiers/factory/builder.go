package factory

import (
	"github.com/SinfulDragon/dice-calculator/internal/core/common"
)

type Builder interface {
	Build(args ModifierArgs) (common.Modifier, error)
}

type ModifierArgs struct {
	Positional []any
	Named      map[string]any
}
