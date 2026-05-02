# Dice Calculator

A multi-functional dice rolling tool with an extensible modifier system and statistical analysis.

## Description

The program allows rolling dice of various types, applying modifiers (rerolls, etc.), combining results via arithmetic operations, and analyzing probability distributions.

## Current Features

- **Dice rolling** — support for any number of dice with any number of sides (`NdM` format, e.g. `2d12`)
- **Arithmetic operations** — addition, subtraction, multiplication, division of roll results
- **Modifier system** — extensible architecture for applying modifiers to dice
- **Statistics** — probability distribution and value analysis

## Architecture

The project is built around an **expression tree** using the Composite pattern:

```
FormulaNode (interface)
├── DiceNode — dice roll (2d12)
├── BinaryNode — binary arithmetic (+, -, *, /)
├── UnaryNode — unary arithmetic (+, -)
├── FlatNode — fixed numeric value
└── ModifierNode — applies a modifier to a child node
```

### Packages

| Package | Purpose |
|---------|---------|
| `internal/core/common` | Base types (`Die` — a die with sides and value, `Modifier` interface) |
| `internal/core/tree` | Expression tree (nodes and `FormulaNode` interface) |
| `internal/core/modifiers/factory` | Modifier builder system (`Builder` interface, global `BuilderRegistry`) |
| `internal/core/modifiers/reroll` | Reroll modifier implementation (`RerollModifier` with configurable modes) |
| `internal/core/parser` | String formula parser (lexer + Pratt-style parser) |
| `internal/core/stats` | Statistical analysis of rolls (WIP) |
| `internal/core/utils` | Utility helpers for argument parsing |
| `internal/gui` | Graphical user interface (Fyne, placeholder) |

## Reroll Modes

| Mode | Description |
|------|-------------|
| `JustReroll` | Reroll all dice |
| `RerollHighest` | Reroll the highest value |
| `RerollLowest` | Reroll the lowest value |
| `RerollBelow` | Reroll dice with value below a threshold |
| `RerollAbove` | Reroll dice with value above a threshold |
| `RerollExact` | Reroll dice with specific values |

## Usage Example

```go
left := &tree.DiceNode{Raw: "2d12", Count: 2, Sides: 12}
right := &tree.DiceNode{Raw: "1d6", Count: 1, Sides: 6}

formula := &tree.BinaryNode{
    Raw:   "",
    Op:    tree.BinaryPlus,
    Left:  left,
    Right: right,
}

formula.Roll()
fmt.Println(formula.Evaluate()) // 2d12 + 1d6
```

Formula parsing from a string:

```go
tree, err := parser.ParseFormula("2d12 + 1d6.reroll(rerollexact, values:[1,2])")
if err != nil {
    // handle error
}
tree.Roll()
fmt.Println(tree.Evaluate())
```

## Running

```bash
go run cmd/dice-calculator/main.go
```

Run tests:

```bash
go test ./...
```

---

**Language:** Go  
**GUI framework:** Fyne v2

