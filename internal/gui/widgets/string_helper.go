package widgets

import "github.com/SinfulDragon/dice-calculator/internal/core/tree"

func NodeString(n tree.FormulaNode) string {
	if s, ok := n.(interface{ String() string }); ok {
		return s.String()
	}
	return ""
}
