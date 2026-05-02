package parser

import (
	"fmt"
	"strconv"
	"strings"
)

type tokenType int

const (
	tokenEOF tokenType = iota
	tokenNumber
	tokenIdentifier
	tokenDice
	tokenPlus
	tokenMinus
	tokenMul
	tokenDiv
	tokenLParen
	tokenRParen
	tokenLBracket
	tokenRBracket
	tokenComma
	tokenColon
	tokenDot
)

type token struct {
	Type    tokenType
	Literal string
	Value   int
}

type lexer struct {
	input         string
	pos           int
	ch            byte
	stringBuilder strings.Builder
}

func newLexer(input string) *lexer {
	l := lexer{input: input}
	l.readChar()
	l.stringBuilder.Grow(16)
	return &l
}

func (l *lexer) readChar() {
	if l.pos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.pos]
	}
	l.pos++
}

func (l *lexer) NextToken() (token, error) {
	l.skipWhitespace()
	token := token{}

	switch l.ch {
	case 0:
		token.Type = tokenEOF
		token.Literal = ""
	case '+':
		token.Type = tokenPlus
		token.Literal = "+"
	case '-':
		token.Type = tokenMinus
		token.Literal = "-"
	case '*':
		token.Type = tokenMul
		token.Literal = "*"
	case '/':
		token.Type = tokenDiv
		token.Literal = "/"
	case '(':
		token.Type = tokenLParen
		token.Literal = "("
	case ')':
		token.Type = tokenRParen
		token.Literal = ")"
	case '[':
		token.Type = tokenLBracket
		token.Literal = "["
	case ']':
		token.Type = tokenRBracket
		token.Literal = "]"
	case ',':
		token.Type = tokenComma
		token.Literal = ","
	case ':':
		token.Type = tokenColon
		token.Literal = ":"
	case '.':
		token.Type = tokenDot
		token.Literal = "."
	default:
		if l.ch >= '0' && l.ch <= '9' {
			var err error
			token.Literal = l.readNumber()
			token.Type = tokenNumber
			token.Value, err = strconv.Atoi(token.Literal)
			if err != nil {
				return token, err
			}
			return token, nil
		} else if l.ch >= 'a' && l.ch <= 'z' || l.ch >= 'A' && l.ch <= 'Z' {
			token.Literal = l.readIdentifier()
			if token.Literal == "d" || token.Literal == "D" {
				token.Type = tokenDice
			} else {
				token.Type = tokenIdentifier
			}
			return token, nil
		} else {
			return token, fmt.Errorf("unexpected character: %v", l.ch)
		}
	}

	if l.ch != 0 {
		l.readChar()
	}

	return token, nil
}

func (l *lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *lexer) readNumber() string {
	l.stringBuilder.Reset()
	for l.ch >= '0' && l.ch <= '9' {
		l.stringBuilder.WriteByte(l.ch)
		l.readChar()
	}
	return l.stringBuilder.String()
}

func (l *lexer) readIdentifier() string {
	l.stringBuilder.Reset()
	for l.ch >= 'a' && l.ch <= 'z' || l.ch >= 'A' && l.ch <= 'Z' {
		l.stringBuilder.WriteByte(l.ch)
		l.readChar()
	}
	return l.stringBuilder.String()
}
