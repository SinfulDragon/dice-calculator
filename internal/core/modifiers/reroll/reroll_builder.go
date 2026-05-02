package reroll

import (
	"fmt"
	"slices"

	"github.com/SinfulDragon/dice-calculator/internal/core/common"
	"github.com/SinfulDragon/dice-calculator/internal/core/modifiers/factory"
	"github.com/SinfulDragon/dice-calculator/internal/core/utils"
)

var modeNames = map[string]RerollMode{
	"justreroll":    JustReroll,
	"rerollhighest": RerollHighest,
	"rerolllowest":  RerollLowest,
	"rerollbelow":   RerollBelow,
	"rerollabove":   RerollAbove,
	"rerollexact":   RerollExact,
}

var requiredArgs = map[RerollMode][]string{
	RerollBelow: {"minvalue"},
	RerollAbove: {"minvalue"},
	RerollExact: {"values"},
}

type RerollBuilder struct{}

func (b *RerollBuilder) Build(args factory.ModifierArgs) (common.Modifier, error) {
	modifier := RerollModifier{}

	// parse positional arguments
	if len(args.Positional) == 0 {
		return nil, fmt.Errorf("no reroll mode specified")
	}
	modeStr, ok := args.Positional[0].(string)
	if !ok {
		return nil, fmt.Errorf("reroll mode must be a string, got %T", args.Positional[0])
	}
	rerollMode, ok := modeNames[modeStr]
	if !ok {
		return nil, fmt.Errorf("unknown reroll mode: %s", args.Positional[0])
	}
	modifier.Mode = rerollMode

	// parse named arguments
	if args.Named != nil {
		var err error
		modifier.Limit, err = utils.ArgInt(args.Named, "limit", slices.Contains(requiredArgs[rerollMode], "limit"))
		if err != nil {
			return nil, err
		}
		modifier.MinValue, err = utils.ArgInt(args.Named, "minvalue", slices.Contains(requiredArgs[rerollMode], "minvalue"))
		if err != nil {
			return nil, err
		}
		modifier.Values, err = utils.ArgIntSlice(args.Named, "values", slices.Contains(requiredArgs[rerollMode], "values"))
		if err != nil {
			return nil, err
		}
	}

	return &modifier, nil
}

func (b *RerollBuilder) Schema() factory.ModifierSchema {
	return factory.ModifierSchema{
		Name:        "reroll",
		Description: "Reroll dice based on value conditions",
		Args: []factory.ArgSchema{
			{
				Name:     "mode",
				Type:     factory.ArgEnum,
				Required: true,
				Options:  []string{"justreroll", "rerollhighest", "rerolllowest", "rerollbelow", "rerollabove", "rerollexact"},
			},
			{
				Name:     "minvalue",
				Type:     factory.ArgInt,
				Required: false,
			},
			{
				Name:     "values",
				Type:     factory.ArgIntSlice,
				Required: false,
			},
			{
				Name:     "limit",
				Type:     factory.ArgInt,
				Required: false,
			},
		},
	}
}
