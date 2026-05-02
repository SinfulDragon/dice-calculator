package factory

import (
	"github.com/SinfulDragon/dice-calculator/internal/core/common"
)

type ArgType string

const (
	ArgInt      ArgType = "int"
	ArgIntSlice ArgType = "[]int"
	ArgString   ArgType = "string"
	ArgEnum     ArgType = "enum"
)

type Builder interface {
	Build(args ModifierArgs) (common.Modifier, error)
	Schema() ModifierSchema
}

type ArgSchema struct {
	Name     string
	Type     ArgType
	Required bool
	Options  []string
}

type ModifierSchema struct {
	Name        string
	Description string
	Args        []ArgSchema
}

type ModifierArgs struct {
	Positional []any
	Named      map[string]any
}
