package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SinfulDragon/dice-calculator/internal/core/modifiers/factory"
	"github.com/SinfulDragon/dice-calculator/internal/core/modifiers/reroll"
	"github.com/SinfulDragon/dice-calculator/internal/core/tree"
)

func init() {
	factory.GlobalRegistry.Register("reroll", &reroll.RerollBuilder{})
}

// assertNodeEqual рекурсивно сравнивает два AST-дерева.
func assertNodeEqual(t *testing.T, expected, actual tree.FormulaNode) {
	t.Helper()
	if expected == nil && actual == nil {
		return
	}
	require.NotNil(t, actual, "expected non-nil node")

	switch exp := expected.(type) {
	case *tree.FlatNode:
		act, ok := actual.(*tree.FlatNode)
		require.True(t, ok, "expected *FlatNode, got %T", actual)
		assert.Equal(t, exp.Value, act.Value)

	case *tree.DiceNode:
		act, ok := actual.(*tree.DiceNode)
		require.True(t, ok, "expected *DiceNode, got %T", actual)
		assert.Equal(t, exp.Count, act.Count)
		assert.Equal(t, exp.Sides, act.Sides)

	case *tree.UnaryNode:
		act, ok := actual.(*tree.UnaryNode)
		require.True(t, ok, "expected *UnaryNode, got %T", actual)
		assert.Equal(t, exp.Op, act.Op)
		assertNodeEqual(t, exp.Child, act.Child)

	case *tree.BinaryNode:
		act, ok := actual.(*tree.BinaryNode)
		require.True(t, ok, "expected *BinaryNode, got %T", actual)
		assert.Equal(t, exp.Op, act.Op)
		assertNodeEqual(t, exp.Left, act.Left)
		assertNodeEqual(t, exp.Right, act.Right)

	case *tree.ModifierNode:
		act, ok := actual.(*tree.ModifierNode)
		require.True(t, ok, "expected *ModifierNode, got %T", actual)
		assertNodeEqual(t, exp.Child, act.Child)
		// Проверяем, что модификатор того же Go-типа
		assert.IsType(t, exp.Modifier, act.Modifier)

	default:
		t.Fatalf("unknown node type in assertion: %T", expected)
	}
}

func TestParseFormula_Success(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected tree.FormulaNode
	}{
		{
			name:     "1. Простое число",
			input:    "42",
			expected: &tree.FlatNode{Value: 42},
		},
		{
			name:     "2. Неявное количество дайсов",
			input:    "d6",
			expected: &tree.DiceNode{Count: 1, Sides: 6},
		},
		{
			name:     "3. Явное количество и грани",
			input:    "2d12",
			expected: &tree.DiceNode{Count: 2, Sides: 12},
		},
		{
			name:  "4. Сложение",
			input: "2d12 + 4",
			expected: &tree.BinaryNode{
				Op:    tree.BinaryPlus,
				Left:  &tree.DiceNode{Count: 2, Sides: 12},
				Right: &tree.FlatNode{Value: 4},
			},
		},
		{
			name:  "5. Умножение",
			input: "3 * d6",
			expected: &tree.BinaryNode{
				Op:    tree.BinaryMul,
				Left:  &tree.FlatNode{Value: 3},
				Right: &tree.DiceNode{Count: 1, Sides: 6},
			},
		},
		{
			name:  "6. Левоассоциативность минуса",
			input: "10 - 2 - 3",
			expected: &tree.BinaryNode{
				Op:    tree.BinaryMinus,
				Left:  &tree.BinaryNode{Op: tree.BinaryMinus, Left: &tree.FlatNode{Value: 10}, Right: &tree.FlatNode{Value: 2}},
				Right: &tree.FlatNode{Value: 3},
			},
		},
		{
			name:  "7. Приоритет умножения над сложением",
			input: "2 * 3 + 4",
			expected: &tree.BinaryNode{
				Op:    tree.BinaryPlus,
				Left:  &tree.BinaryNode{Op: tree.BinaryMul, Left: &tree.FlatNode{Value: 2}, Right: &tree.FlatNode{Value: 3}},
				Right: &tree.FlatNode{Value: 4},
			},
		},
		{
			name:  "8. Унарный минус перед дайсом",
			input: "-1d6",
			expected: &tree.UnaryNode{
				Op:    tree.UnaryMinus,
				Child: &tree.DiceNode{Count: 1, Sides: 6},
			},
		},
		{
			name:  "9. Двойной унарный минус",
			input: "--2",
			expected: &tree.UnaryNode{
				Op: tree.UnaryMinus,
				Child: &tree.UnaryNode{
					Op:    tree.UnaryMinus,
					Child: &tree.FlatNode{Value: 2},
				},
			},
		},
		{
			name:  "10. Скобки изменяют приоритет",
			input: "(2d6 + 1) * 3",
			expected: &tree.BinaryNode{
				Op: tree.BinaryMul,
				Left: &tree.BinaryNode{
					Op:    tree.BinaryPlus,
					Left:  &tree.DiceNode{Count: 2, Sides: 6},
					Right: &tree.FlatNode{Value: 1},
				},
				Right: &tree.FlatNode{Value: 3},
			},
		},
		{
			name:  "11. Цепочка сложения трёх элементов",
			input: "2d12 + 1d6 + 4",
			expected: &tree.BinaryNode{
				Op:    tree.BinaryPlus,
				Left:  &tree.BinaryNode{Op: tree.BinaryPlus, Left: &tree.DiceNode{Count: 2, Sides: 12}, Right: &tree.DiceNode{Count: 1, Sides: 6}},
				Right: &tree.FlatNode{Value: 4},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := ParseFormula(tt.input)
			require.NoError(t, err, "parse error: %v", err)
			assertNodeEqual(t, tt.expected, node)
		})
	}
}

func TestParseFormula_Modifiers(t *testing.T) {
	// 12. Модификатор reroll с аргументами
	t.Run("reroll with args", func(t *testing.T) {
		mod, err := factory.GlobalRegistry.Build("reroll", factory.ModifierArgs{
			Positional: []any{"rerollexact"},
			Named:      map[string]any{"values": []any{1, 2}},
		})
		require.NoError(t, err)
		expected := &tree.ModifierNode{
			Modifier: mod,
			Child:    &tree.DiceNode{Count: 1, Sides: 6},
		}
		node, err := ParseFormula("1d6.reroll(RerollExact, Values:[1, 2])")
		require.NoError(t, err)
		assertNodeEqual(t, expected, node)
	})

	// 13. Цепочка модификаторов
	t.Run("modifier chain", func(t *testing.T) {
		innerMod, err := factory.GlobalRegistry.Build("reroll", factory.ModifierArgs{
			Positional: []any{"justreroll"},
		})
		require.NoError(t, err)
		outerMod, err := factory.GlobalRegistry.Build("reroll", factory.ModifierArgs{
			Positional: []any{"justreroll"},
		})
		require.NoError(t, err)
		inner := &tree.ModifierNode{Modifier: innerMod, Child: &tree.DiceNode{Count: 1, Sides: 6}}
		expected := &tree.ModifierNode{Modifier: outerMod, Child: inner}
		node, err := ParseFormula("1d6.reroll(JustReroll).reroll(JustReroll)")
		require.NoError(t, err)
		assertNodeEqual(t, expected, node)
	})
}

func TestParseFormula_Errors(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedErr string
	}{
		{
			name:        "1. Пустая строка",
			input:       "",
			expectedErr: "expected number or `(`, got",
		},
		{
			name:        "2. Дайс без количества граней",
			input:       "2d",
			expectedErr: "expected number after `d`",
		},
		{
			name:        "3. Двойная точка",
			input:       "2d6..reroll()",
			expectedErr: "expected identifier after `.`",
		},
		{
			name:        "4. Незакрытая скобка модификатора",
			input:       "1d6.reroll(",
			expectedErr: "failed to parse modifier value: unexpected token:",
		},
		{
			name:        "5. Пустое значение named аргумента",
			input:       "1d6.reroll(values:)",
			expectedErr: "unexpected token",
		},
		{
			name:        "6. Незакрытая скобка выражения",
			input:       "(2d6",
			expectedErr: "expected `)` after expression",
		},
		{
			name:        "7. Лишний токен после выражения",
			input:       "42 7",
			expectedErr: "unexpected trailing token",
		},
		{
			name:        "8. Пустой reroll без обязательного mode",
			input:       "1d6.reroll()",
			expectedErr: "no reroll mode specified",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := ParseFormula(tt.input)
			assert.Nil(t, node)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestParseFormula_UnaryPlus(t *testing.T) {
	node, err := ParseFormula("+3")
	require.NoError(t, err)
	assertNodeEqual(t, &tree.UnaryNode{Op: tree.UnaryPlus, Child: &tree.FlatNode{Value: 3}}, node)
}

func TestParseFormula_WhitespaceTolerance(t *testing.T) {
	node, err := ParseFormula("  2  d 12  +  4  ")
	require.NoError(t, err)
	assertNodeEqual(t, &tree.BinaryNode{
		Op:    tree.BinaryPlus,
		Left:  &tree.DiceNode{Count: 2, Sides: 12},
		Right: &tree.FlatNode{Value: 4},
	}, node)
}

func TestParseFormula_Division(t *testing.T) {
	node, err := ParseFormula("8 / 2")
	require.NoError(t, err)
	assertNodeEqual(t, &tree.BinaryNode{
		Op:    tree.BinaryDiv,
		Left:  &tree.FlatNode{Value: 8},
		Right: &tree.FlatNode{Value: 2},
	}, node)
}

func TestParseFormula_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected tree.FormulaNode
	}{
		{
			name:     "0d6 — ноль дайсов",
			input:    "0d6",
			expected: &tree.DiceNode{Count: 0, Sides: 6},
		},
		{
			name:     "d0 — ноль граней",
			input:    "d0",
			expected: &tree.DiceNode{Count: 1, Sides: 0},
		},
		{
			name:     "1d1 — минимальный дайс",
			input:    "1d1",
			expected: &tree.DiceNode{Count: 1, Sides: 1},
		},
		{
			name:  "+d6 — унарный плюс перед dice",
			input: "+d6",
			expected: &tree.UnaryNode{
				Op:    tree.UnaryPlus,
				Child: &tree.DiceNode{Count: 1, Sides: 6},
			},
		},
		{
			name:  "-(3+4) — унарный минус перед скобками",
			input: "-(3+4)",
			expected: &tree.UnaryNode{
				Op: tree.UnaryMinus,
				Child: &tree.BinaryNode{
					Op:    tree.BinaryPlus,
					Left:  &tree.FlatNode{Value: 3},
					Right: &tree.FlatNode{Value: 4},
				},
			},
		},
		{
			name:  "-(-2) — двойной унарный минус в скобках",
			input: "-(-2)",
			expected: &tree.UnaryNode{
				Op: tree.UnaryMinus,
				Child: &tree.UnaryNode{
					Op:    tree.UnaryMinus,
					Child: &tree.FlatNode{Value: 2},
				},
			},
		},
		{
			name:  "8/4*2 — левоассоциативность деления и умножения",
			input: "8/4*2",
			expected: &tree.BinaryNode{
				Op:   tree.BinaryMul,
				Left: &tree.BinaryNode{Op: tree.BinaryDiv, Left: &tree.FlatNode{Value: 8}, Right: &tree.FlatNode{Value: 4}},
				Right: &tree.FlatNode{Value: 2},
			},
		},
		{
			name:  "8-4+2 — левоассоциативность минуса и плюса",
			input: "8-4+2",
			expected: &tree.BinaryNode{
				Op:   tree.BinaryPlus,
				Left: &tree.BinaryNode{Op: tree.BinaryMinus, Left: &tree.FlatNode{Value: 8}, Right: &tree.FlatNode{Value: 4}},
				Right: &tree.FlatNode{Value: 2},
			},
		},
		{
			name:  "3+(2*4) — скобки внутри выражения",
			input: "3+(2*4)",
			expected: &tree.BinaryNode{
				Op:    tree.BinaryPlus,
				Left:  &tree.FlatNode{Value: 3},
				Right: &tree.BinaryNode{Op: tree.BinaryMul, Left: &tree.FlatNode{Value: 2}, Right: &tree.FlatNode{Value: 4}},
			},
		},
		{
			name:     "(((42))) — тройные скобки",
			input:    "(((42)))",
			expected: &tree.FlatNode{Value: 42},
		},
		{
			name:  "-+-3 — чередование унарных операторов",
			input: "-+-3",
			expected: &tree.UnaryNode{
				Op: tree.UnaryMinus,
				Child: &tree.UnaryNode{
					Op: tree.UnaryPlus,
					Child: &tree.UnaryNode{
						Op:    tree.UnaryMinus,
						Child: &tree.FlatNode{Value: 3},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := ParseFormula(tt.input)
			require.NoError(t, err, "parse error: %v", err)
			assertNodeEqual(t, tt.expected, node)
		})
	}
}

func TestParseFormula_EdgeCaseErrors(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedErr string
	}{
		{
			name:        "Пустые скобки",
			input:       "()",
			expectedErr: "expected number or `(`, got",
		},
		{
			name:        "Лишняя правая скобка",
			input:       "2d6)",
			expectedErr: "unexpected trailing token",
		},
		{
			name:        "Бинарный оператор без правого операнда",
			input:       "2d6 +",
			expectedErr: "expected number or `(`, got",
		},
		{
			name:        "Два бинарных оператора подряд",
			input:       "2d6 + * 3",
			expectedErr: "expected number or `(`, got",
		},
		{
			name:        "Модификатор без скобок",
			input:       "1d6.reroll",
			expectedErr: "expected `(` after identifier",
		},
		{
			name:        "Запятая в конце аргументов",
			input:       "1d6.reroll(1,)",
			expectedErr: "unexpected token",
		},
		{
			name:        "Незакрытый массив",
			input:       "1d6.reroll([1,2)",
			expectedErr: "expected ',' or ']' in array",
		},
		{
			name:        "Пустая запятая в аргументах",
			input:       "1d6.reroll(,)",
			expectedErr: "unexpected token",
		},
		{
			name:        "Точка после flat числа (не дайс)",
			input:       "5.reroll()",
			expectedErr: "unexpected trailing token",
		},
		{
			name:        "Неподдерживаемый символ в формуле",
			input:       "2d6 $ 3",
			expectedErr: "unexpected character",
		},
		{
			name:        "d без числа граней",
			input:       "d",
			expectedErr: "expected number after `d`",
		},
		{
			name:        "Два d подряд без числа между ними",
			input:       "2d6d8",
			expectedErr: "unexpected trailing token",
		},
		{
			name:        "Неизвестный модификатор",
			input:       "1d6.unknown()",
			expectedErr: "unknown modifier: unknown",
		},
		{
			name:        "Две лишние правые скобки",
			input:       "(2d6))",
			expectedErr: "unexpected trailing token",
		},
		{
			name:        "Пропущенная запятая в массиве",
			input:       "1d6.reroll(RerollExact,Values:[1 2])",
			expectedErr: "expected ',' or ']' in array",
		},
		{
			name:        "Пустой массив как positional arg",
			input:       "1d6.reroll([])",
			expectedErr: "reroll mode must be a string",
		},
		{
			name:        "Число вместо массива для named arg values",
			input:       "1d6.reroll(RerollExact,Values:1)",
			expectedErr: "parameter 'values' must be a slice of numbers",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := ParseFormula(tt.input)
			assert.Nil(t, node)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}
