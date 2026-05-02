package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertTokens(t *testing.T, input string, expected []token) {
	l := newLexer(input)
	for i, exp := range expected {
		token, err := l.NextToken()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		assert.Equal(t, exp.Type, token.Type, "token %d type mismatch", i)
		assert.Equal(t, exp.Literal, token.Literal, "token %d literal mismatch", i)
		if exp.Type == tokenNumber {
			assert.Equal(t, exp.Value, token.Value, "token %d value mismatch", i)
		}
	}
	final, err := l.NextToken()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	assert.Equal(t, tokenEOF, final.Type, "expected EOF")
}

func TestLexer_NextToken(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []token
	}{{
		name:  "1. Простое число",
		input: "42",
		expected: []token{
			{Type: tokenNumber, Literal: "42", Value: 42},
		},
	},
		{
			name:  "2. Простой дайс",
			input: "d6",
			expected: []token{
				{Type: tokenDice, Literal: "d"},
				{Type: tokenNumber, Literal: "6", Value: 6},
			},
		},
		{
			name:  "3. Дайс с количеством",
			input: "2d12",
			expected: []token{
				{Type: tokenNumber, Literal: "2", Value: 2},
				{Type: tokenDice, Literal: "d"},
				{Type: tokenNumber, Literal: "12", Value: 12},
			},
		},
		{
			name:  "4. Арифметика",
			input: "2d12 + 4",
			expected: []token{
				{Type: tokenNumber, Literal: "2", Value: 2},
				{Type: tokenDice, Literal: "d"},
				{Type: tokenNumber, Literal: "12", Value: 12},
				{Type: tokenPlus, Literal: "+"},
				{Type: tokenNumber, Literal: "4", Value: 4},
			},
		},
		{
			name:  "5. Произведение",
			input: "3 * d6",
			expected: []token{
				{Type: tokenNumber, Literal: "3", Value: 3},
				{Type: tokenMul, Literal: "*"},
				{Type: tokenDice, Literal: "d"},
				{Type: tokenNumber, Literal: "6", Value: 6},
			},
		},
		{
			name:  "6. Унарный минус",
			input: "-1d6",
			expected: []token{
				{Type: tokenMinus, Literal: "-"},
				{Type: tokenNumber, Literal: "1", Value: 1},
				{Type: tokenDice, Literal: "d"},
				{Type: tokenNumber, Literal: "6", Value: 6},
			},
		},
		{
			name:  "7. Скобки",
			input: "(2d6)",
			expected: []token{
				{Type: tokenLParen, Literal: "("},
				{Type: tokenNumber, Literal: "2", Value: 2},
				{Type: tokenDice, Literal: "d"},
				{Type: tokenNumber, Literal: "6", Value: 6},
				{Type: tokenRParen, Literal: ")"},
			},
		},
		{
			name:  "8. Модификатор (полный)",
			input: "reroll(RerollExact, Values:[1, 2, 3])",
			expected: []token{
				{Type: tokenIdentifier, Literal: "reroll"},
				{Type: tokenLParen, Literal: "("},
				{Type: tokenIdentifier, Literal: "RerollExact"},
				{Type: tokenComma, Literal: ","},
				{Type: tokenIdentifier, Literal: "Values"},
				{Type: tokenColon, Literal: ":"},
				{Type: tokenLBracket, Literal: "["},
				{Type: tokenNumber, Literal: "1", Value: 1},
				{Type: tokenComma, Literal: ","},
				{Type: tokenNumber, Literal: "2", Value: 2},
				{Type: tokenComma, Literal: ","},
				{Type: tokenNumber, Literal: "3", Value: 3},
				{Type: tokenRBracket, Literal: "]"},
				{Type: tokenRParen, Literal: ")"},
			},
		},
		{
			name:  "9. Цепочка",
			input: "1d6.reroll()",
			expected: []token{
				{Type: tokenNumber, Literal: "1", Value: 1},
				{Type: tokenDice, Literal: "d"},
				{Type: tokenNumber, Literal: "6", Value: 6},
				{Type: tokenDot, Literal: "."},
				{Type: tokenIdentifier, Literal: "reroll"},
				{Type: tokenLParen, Literal: "("},
				{Type: tokenRParen, Literal: ")"},
			},
		},
		{
			name:     "10. Пустая строка",
			input:    "",
			expected: []token{},
		},
		{
			name:  "11. Пробелы",
			input: "  2   d   12  ",
			expected: []token{
				{Type: tokenNumber, Literal: "2", Value: 2},
				{Type: tokenDice, Literal: "d"},
				{Type: tokenNumber, Literal: "12", Value: 12},
			},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertTokens(t, tt.input, tt.expected)
		})
	}
}

func TestLexer_NextToken_Errors(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		errMsg  string
		tokenAt int // на каком по счету токене ожидаем ошибку (0-based)
	}{
		{
			name:    "Unicode символ",
			input:   "2d6+💀",
			errMsg:  "unexpected character",
			tokenAt: 4, // 2 d 6 + 💀
		},
		{
			name:    "Неподдерживаемый спецсимвол @",
			input:   "2d6@3",
			errMsg:  "unexpected character",
			tokenAt: 3,
		},
		{
			name:    "Vertical tab не считается пробелом",
			input:   "2\vd6",
			errMsg:  "unexpected character",
			tokenAt: 1,
		},
		{
			name:    "Form feed не считается пробелом",
			input:   "2\fd6",
			errMsg:  "unexpected character",
			tokenAt: 1,
		},
		{
			name:    "Underscore в идентификаторе",
			input:   "foo_bar",
			errMsg:  "unexpected character",
			tokenAt: 1,
		},
		{
			name:    "Переполнение int",
			input:   "999999999999999999999",
			errMsg:  "value out of range",
			tokenAt: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := newLexer(tt.input)
			for i := 0; i < tt.tokenAt; i++ {
				_, err := l.NextToken()
				require.NoError(t, err, "unexpected error before target token at %d", i)
			}
			_, err := l.NextToken()
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}

func TestLexer_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []token
	}{
		{
			name:  "Ведущие нули",
			input: "007",
			expected: []token{
				{Type: tokenNumber, Literal: "007", Value: 7},
			},
		},
		{
			name:  "Ноль",
			input: "0",
			expected: []token{
				{Type: tokenNumber, Literal: "0", Value: 0},
			},
		},
		{
			name:  "Точка в начале строки",
			input: ".5",
			expected: []token{
				{Type: tokenDot, Literal: "."},
				{Type: tokenNumber, Literal: "5", Value: 5},
			},
		},
		{
			name:  "Три точки подряд",
			input: "...",
			expected: []token{
				{Type: tokenDot, Literal: "."},
				{Type: tokenDot, Literal: "."},
				{Type: tokenDot, Literal: "."},
			},
		},
		{
			name:  "Смешанный регистр идентификатора",
			input: "ReRoLl",
			expected: []token{
				{Type: tokenIdentifier, Literal: "ReRoLl"},
			},
		},
		{
			name:  "Число, примыкающее к букве (не d)",
			input: "123abc",
			expected: []token{
				{Type: tokenNumber, Literal: "123", Value: 123},
				{Type: tokenIdentifier, Literal: "abc"},
			},
		},
		{
			name:  "Буква, примыкающая к числу",
			input: "abc123",
			expected: []token{
				{Type: tokenIdentifier, Literal: "abc"},
				{Type: tokenNumber, Literal: "123", Value: 123},
			},
		},
		{
			name:  "Большое, но валидное число",
			input: "2147483647",
			expected: []token{
				{Type: tokenNumber, Literal: "2147483647", Value: 2147483647},
			},
		},
		{
			name:     "Только пробельные символы",
			input:    "   \t\n\r  ",
			expected: []token{},
		},
		{
			name:  "Каретка (CR) как пробел",
			input: "2\r+\r3",
			expected: []token{
				{Type: tokenNumber, Literal: "2", Value: 2},
				{Type: tokenPlus, Literal: "+"},
				{Type: tokenNumber, Literal: "3", Value: 3},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertTokens(t, tt.input, tt.expected)
		})
	}
}
