package token

import (
	"fmt"
	"go-compiler/main/errors"
)

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

func NewToken(t TokenType, lexeme string, literal interface{}, line int) *Token {
	return &Token{Type: t, Lexeme: lexeme, Literal: literal, Line: line}
}

func (t *Token) String() string {
	return fmt.Sprintf("%v %v %v", t.Type, t.Lexeme, t.Literal)
}

func TokenError(t *Token, message string) string {
	if t.Type == EOF {
		return errors.Report(t.Line, " at end", message)
	} else {
		return errors.Report(t.Line, " at '"+t.Lexeme+"'", message)
	}
}
