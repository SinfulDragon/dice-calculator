package tree

import (
	"fmt"
	"strings"

	"github.com/SinfulDragon/dice-calculator/internal/core/common"
	"github.com/SinfulDragon/dice-calculator/internal/core/modifiers/factory"
)

type ModifierNode struct {
	Raw      string
	Modifier common.Modifier
	Child    FormulaNode
	Name     string
	Args     factory.ModifierArgs
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
		Name:     n.Name,
		Args:     n.Args,
	}
}

func formatModifierArgs(args factory.ModifierArgs) string {
	var parts []string
	for _, p := range args.Positional {
		switch v := p.(type) {
		case string:
			parts = append(parts, v)
		case int:
			parts = append(parts, fmt.Sprintf("%d", v))
		case []any:
			var inner []string
			for _, iv := range v {
				switch ivv := iv.(type) {
				case int:
					inner = append(inner, fmt.Sprintf("%d", ivv))
				default:
					inner = append(inner, fmt.Sprintf("%v", ivv))
				}
			}
			parts = append(parts, fmt.Sprintf("[%s]", strings.Join(inner, ",")))
		default:
			parts = append(parts, fmt.Sprintf("%v", p))
		}
	}
	for k, v := range args.Named {
		switch val := v.(type) {
		case string:
			parts = append(parts, fmt.Sprintf("%s:%s", k, val))
		case int:
			parts = append(parts, fmt.Sprintf("%s:%d", k, val))
		case []any:
			var inner []string
			for _, iv := range val {
				switch ivv := iv.(type) {
				case int:
					inner = append(inner, fmt.Sprintf("%d", ivv))
				default:
					inner = append(inner, fmt.Sprintf("%v", ivv))
				}
			}
			parts = append(parts, fmt.Sprintf("%s:[%s]", k, strings.Join(inner, ",")))
		default:
			parts = append(parts, fmt.Sprintf("%s:%v", k, val))
		}
	}
	return strings.Join(parts, ", ")
}

func (n *ModifierNode) String() string {
	if n.Name == "" {
		return stringNode(n.Child)
	}
	return fmt.Sprintf("%s.%s(%s)", stringNode(n.Child), n.Name, formatModifierArgs(n.Args))
}
