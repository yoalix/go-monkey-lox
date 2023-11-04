package scanner

import (
	"strconv"

	"go-compiler/main/errors"
	"go-compiler/main/token"
)

type Scanner struct {
	source  string
	tokens  []*token.Token
	start   int
	current int
	line    int
}

const ()

func NewScanner(source string) *Scanner {
	return &Scanner{source: source, tokens: []*token.Token{}, start: 0, current: 0, line: 1}
}

func (s *Scanner) ScanTokens() []*token.Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.ScanToken()
	}
	s.tokens = append(s.tokens, token.NewToken(token.EOF, "", nil, s.line))
	return s.tokens
}

func (s *Scanner) ScanToken() *token.Token {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(token.LEFT_PAREN)
	case ')':
		s.addToken(token.RIGHT_PAREN)
	case '{':
		s.addToken(token.LEFT_BRACE)
	case '}':
		s.addToken(token.RIGHT_BRACE)
	case '[':
		s.addToken(token.LEFT_BRACKET)
	case ']':
		s.addToken(token.RIGHT_BRACKET)
	case ',':
		s.addToken(token.COMMA)
	case '.':
		s.addToken(token.DOT)
	case '-':
		s.addToken(token.MINUS)
	case '+':
		s.addToken(token.PLUS)
	case ';':
		s.addToken(token.SEMICOLON)
	case ':':
		s.addToken(token.COLON)
	case '*':
		s.addToken(token.STAR)
	case '!':
		if s.match('=') {
			s.addToken(token.BANG_EQUAL)
		} else {
			s.addToken(token.BANG)
		}
	case '=':
		if s.match('=') {
			s.addToken(token.EQUAL_EQUAL)
		} else {
			s.addToken(token.EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addToken(token.LESS_EQUAL)
		} else {
			s.addToken(token.LESS)
		}
	case '>':
		if s.match('=') {
			s.addToken(token.GREATER_EQUAL)
		} else {
			s.addToken(token.GREATER)
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(token.SLASH)
		}
	case ' ', '\r', '\t':
		// Ignore whitespace
		return nil
	case '\n':
		s.line++
		return nil
	case '"':
		s.string()
	default:
		if isDigit(c) {
			s.number()
		} else if isAlpha(c) {
			s.identifier()
		} else {
			errors.Error(s.line, "Unexpected character.")
			return nil
		}
	}
	return s.tokens[len(s.tokens)-1]
}

func (s *Scanner) advance() byte {
	char := s.source[s.current]
	s.current++
	return char
}

func (s *Scanner) addToken(t token.TokenType) {
	s.addTokenLiteral(t, nil)
}
func (s *Scanner) addTokenLiteral(t token.TokenType, literal interface{}) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, token.NewToken(t, text, literal, s.line))
}
func (s *Scanner) addStringLiteral(t token.TokenType, literal string) {
	s.tokens = append(s.tokens, token.NewToken(t, literal, literal, s.line))
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() || s.source[s.current] != expected {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) isAtEnd() bool {
	return len(s.source) <= s.current
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return '\x00'
	}
	return s.source[s.current]
}
func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return '\x00'
	}
	return s.source[s.current+1]
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		errors.Error(s.line, "Unterminated string.")
		return
	}

	// The closing quote
	s.advance()

	// Trim the surrounding quotes
	value := s.source[s.start+1 : s.current-1]
	s.addStringLiteral(token.STRING, value)
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) number() {
	for isDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()
		for isDigit(s.peek()) {
			s.advance()
		}
	}
	floatNumber, err := strconv.ParseFloat(s.source[s.start:s.current], 64)

	if err != nil {
		errors.Error(s.line, "Unable to parse number.")
		return
	}
	s.addTokenLiteral(token.NUMBER, floatNumber)
}

func (s *Scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	tokenType, ok := token.Keywords[text]
	if !ok {
		tokenType = token.IDENTIFIER
	}
	s.addToken(tokenType)

}

func isAlpha(c byte) bool {
	return c >= 'a' && c <= 'z' ||
		c >= 'A' && c <= 'Z' ||
		c == '_'
}

func isAlphaNumeric(c byte) bool {
	return isAlpha(c) || isDigit(c)
}
