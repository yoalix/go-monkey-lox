package token

type TokenType string

const (
	// Single-character tokens.
	LEFT_PAREN    TokenType = "("
	RIGHT_PAREN   TokenType = ")"
	LEFT_BRACE    TokenType = "{"
	RIGHT_BRACE   TokenType = "}"
	LEFT_BRACKET  TokenType = "["
	RIGHT_BRACKET TokenType = "]"
	COMMA         TokenType = ","
	DOT           TokenType = "."
	MINUS         TokenType = "-"
	PLUS          TokenType = "+"
	SEMICOLON     TokenType = ";"
	SLASH         TokenType = "/"
	STAR          TokenType = "*"
	COLON         TokenType = ":"

	// One or two character tokens.
	BANG          TokenType = "!"
	BANG_EQUAL    TokenType = "!="
	EQUAL         TokenType = "="
	EQUAL_EQUAL   TokenType = "=="
	GREATER       TokenType = ">"
	GREATER_EQUAL TokenType = ">="
	LESS          TokenType = "<"
	LESS_EQUAL    TokenType = "<="

	// Literals.
	IDENTIFIER TokenType = "IDENTIFIER"
	STRING     TokenType = "STRING"
	NUMBER     TokenType = "NUMBER"

	// Keywords.
	AND      TokenType = "AND"
	STRUCT   TokenType = "STRUCT"
	ELSE     TokenType = "ELSE"
	FALSE    TokenType = "FALSE"
	FUNCTION TokenType = "FN"
	FOR      TokenType = "FOR"
	IF       TokenType = "IF"
	NIL      TokenType = "NIL"
	OR       TokenType = "OR"
	PRINT    TokenType = "PRINT"
	RETURN   TokenType = "RETURN"
	SUPER    TokenType = "SUPER"
	THIS     TokenType = "THIS"
	TRUE     TokenType = "TRUE"
	LET      TokenType = "LET"
	WHILE    TokenType = "WHILE"

	EOF     TokenType = "EOF"
	ILLEGAL TokenType = "ILLEGAL"
)

var Keywords = map[string]TokenType{
	"and":    AND,
	"struct": STRUCT,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fn":     FUNCTION,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"let":    LET,
}
