package parser

import (
	"fmt"
	"strings"

	"github.com/SinfulDragon/dice-calculator/internal/core/common"
	"github.com/SinfulDragon/dice-calculator/internal/core/modifiers/factory"
	"github.com/SinfulDragon/dice-calculator/internal/core/tree"
)

type Parser struct {
	l         *lexer
	curToken  token
	peekToken token
}

func newParser(l *lexer) (*Parser, error) {
	var err error
	p := &Parser{l: l}
	p.peekToken, err = l.NextToken()
	if err != nil {
		return nil, err
	}
	p.curToken = p.peekToken
	p.peekToken, err = l.NextToken()
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Parser) nextToken() error {
	var err error
	p.curToken = p.peekToken
	p.peekToken, err = p.l.NextToken()
	if err != nil {
		return fmt.Errorf("failed to get next token %w", err)
	}
	return nil
}

// 2d12 + 1d6.reroll(RerollExact, Values:[1, 2, 3]) + 4
func ParseFormula(formula string) (tree.FormulaNode, error) {
	l := newLexer(formula)
	p, err := newParser(l)
	if err != nil {
		return nil, fmt.Errorf("cannot parse formula %s\n%w", formula, err)
	}

	node, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if p.curToken.Type != tokenEOF {
		return nil, fmt.Errorf("unexpected trailing token %s", p.curToken.Literal)
	}

	return node, nil
}

func (p *Parser) parseExpression() (tree.FormulaNode, error) {
	node, err := p.parseAdditive()
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (p *Parser) parseAdditive() (tree.FormulaNode, error) {
	var err error
	left, err := p.parseMultiplicative()
	if err != nil {
		return nil, err
	}
	for p.curToken.Type == tokenPlus || p.curToken.Type == tokenMinus {
		op := p.curToken.Type
		err = p.nextToken()
		if err != nil {
			return nil, err
		}
		right, err := p.parseMultiplicative()
		if err != nil {
			return nil, err
		}
		if op == tokenPlus {
			left = &tree.BinaryNode{Op: tree.BinaryPlus, Left: left, Right: right}
		} else {
			left = &tree.BinaryNode{Op: tree.BinaryMinus, Left: left, Right: right}
		}
	}
	return left, nil
}

func (p *Parser) parseMultiplicative() (tree.FormulaNode, error) {
	var err error
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	for p.curToken.Type == tokenMul || p.curToken.Type == tokenDiv {
		op := p.curToken.Type
		err = p.nextToken()
		if err != nil {
			return nil, err
		}
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		if op == tokenMul {
			left = &tree.BinaryNode{Op: tree.BinaryMul, Left: left, Right: right}
		} else {
			left = &tree.BinaryNode{Op: tree.BinaryDiv, Left: left, Right: right}
		}
	}
	return left, nil
}

func (p *Parser) parseUnary() (tree.FormulaNode, error) {
	var err error
	if p.curToken.Type == tokenPlus || p.curToken.Type == tokenMinus {
		op := p.curToken.Type
		err = p.nextToken()
		if err != nil {
			return nil, err
		}
		child, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		if op == tokenPlus {
			return &tree.UnaryNode{Op: tree.UnaryPlus, Child: child}, nil
		} else {
			return &tree.UnaryNode{Op: tree.UnaryMinus, Child: child}, nil
		}
	}
	return p.parseDiceOrPrimary()
}

func (p *Parser) parseDiceOrPrimary() (tree.FormulaNode, error) {
	// 1d6 and d6 are equivalent
	if p.curToken.Type == tokenNumber && p.peekToken.Type == tokenDice || p.curToken.Type == tokenDice {
		return p.parseDice()
	}
	return p.parsePrimary()
}

func (p *Parser) parseDice() (tree.FormulaNode, error) {
	var err error
	count := 1
	if p.curToken.Type == tokenNumber {
		count = p.curToken.Value
		err = p.nextToken() // eat count
		if err != nil {
			return nil, err
		}
	}
	if p.curToken.Type != tokenDice {
		return nil, fmt.Errorf("expected `d` after dice count, got %s", p.curToken.Literal)
	}
	err = p.nextToken() // eat dice
	if err != nil {
		return nil, err
	}
	if p.curToken.Type != tokenNumber {
		return nil, fmt.Errorf("expected number after `d`, got %s", p.curToken.Literal)
	}
	sides := p.curToken.Value
	err = p.nextToken() // eat sides
	if err != nil {
		return nil, err
	}
	var node tree.FormulaNode
	node = &tree.DiceNode{Count: count, Sides: sides}
	for p.curToken.Type == tokenDot {
		err = p.nextToken() // eat dot
		if err != nil {
			return nil, err
		}
		mod, err := p.parseModifier()
		if err != nil {
			return nil, err
		}
		node = &tree.ModifierNode{Child: node, Modifier: mod}
	}
	return node, nil
}

func (p *Parser) parsePrimary() (tree.FormulaNode, error) {
	var err error
	switch p.curToken.Type {
	case tokenNumber:
		node := &tree.FlatNode{Value: p.curToken.Value}
		err = p.nextToken() // eat number
		if err != nil {
			return nil, err
		}
		return node, nil
	case tokenLParen:
		err = p.nextToken()
		if err != nil {
			return nil, err
		}
		node, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if p.curToken.Type != tokenRParen {
			return nil, fmt.Errorf("expected `)` after expression, got %s", p.curToken.Literal)
		}
		err = p.nextToken() // eat right parenthesis
		if err != nil {
			return nil, err
		}
		return node, nil
	default:
		return nil, fmt.Errorf("expected number or `(`, got %s", p.curToken.Literal)
	}
}

func (p *Parser) parseModifier() (common.Modifier, error) {
	var err error

	if p.curToken.Type != tokenIdentifier {
		return nil, fmt.Errorf("expected identifier after `.`")
	}
	name := strings.ToLower(p.curToken.Literal)
	err = p.nextToken() // eat identifier
	if err != nil {
		return nil, err
	}
	if p.curToken.Type != tokenLParen {
		return nil, fmt.Errorf("expected `(` after identifier")
	}
	err = p.nextToken() // eat left parenthesis
	if err != nil {
		return nil, err
	}
	var args factory.ModifierArgs
	if p.curToken.Type != tokenRParen {
		// parse args
		for {
			if p.curToken.Type == tokenIdentifier && p.peekToken.Type == tokenColon {
				// named args
				key := strings.ToLower(p.curToken.Literal)
				err = p.nextToken() // eat identifier
				if err != nil {
					return nil, err
				}
				err = p.nextToken() // eat colon
				if err != nil {
					return nil, err
				}

				value, err := p.parseModifierValue()
				if err != nil {
					return nil, fmt.Errorf("failed to parse modifier value: %w", err)
				}
				if args.Named == nil {
					args.Named = make(map[string]any)
				}
				args.Named[key] = value
			} else {
				// positional args
				value, err := p.parseModifierValue()
				if err != nil {
					return nil, fmt.Errorf("failed to parse modifier value: %w", err)
				}
				args.Positional = append(args.Positional, value)
			}
			if p.curToken.Type == tokenComma {
				err = p.nextToken() // eat comma
				if err != nil {
					return nil, err
				}
				continue
			}
			break
		}
	}

	if p.curToken.Type != tokenRParen {
		return nil, fmt.Errorf("expected `)` after arguments")
	}
	err = p.nextToken() // eat right parenthesis
	if err != nil {
		return nil, err
	}

	modifier, err := factory.GlobalRegistry.Build(name, args)
	if err != nil {
		return nil, fmt.Errorf("modifier '%s': %w", name, err)
	}
	return modifier, nil
}

func (p *Parser) parseModifierValue() (any, error) {
	var err error
	switch p.curToken.Type {
	case tokenNumber:
		value := p.curToken.Value
		err = p.nextToken() // eat number
		if err != nil {
			return nil, err
		}
		return value, nil
	case tokenIdentifier:
		value := strings.ToLower(p.curToken.Literal)
		err = p.nextToken() // eat identifier
		if err != nil {
			return nil, err
		}
		return value, nil
	case tokenLBracket:
		return p.parseArray()
	default:
		return nil, fmt.Errorf("unexpected token: %s in modifier argument", p.curToken.Literal)
	}
}

func (p *Parser) parseArray() ([]any, error) {
	err := p.nextToken() // eat left bracket
	if err != nil {
		return nil, err
	}
	var elements []any
	for p.curToken.Type != tokenRBracket && p.curToken.Type != tokenEOF {
		element, err := p.parseModifierValue()
		if err != nil {
			return nil, fmt.Errorf("failed to parse array element: %w", err)
		}
		elements = append(elements, element)
		if p.curToken.Type == tokenComma {
			err = p.nextToken() // eat comma
			if err != nil {
				return nil, err
			}
			continue
		}
		if p.curToken.Type != tokenRBracket {
			return nil, fmt.Errorf("expected ',' or ']' in array")
		}
	}
	if p.curToken.Type != tokenRBracket {
		return nil, fmt.Errorf("expected ']' to close array")
	}
	err = p.nextToken() // eat right bracket
	if err != nil {
		return nil, err
	}
	return elements, nil
}
