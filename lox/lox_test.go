package lox

// import (
// 	"go-compiler/main/token"
// 	"testing"
// )

// func TestNextToken(t *testing.T) {
// 	input := `let five = 5;
// let ten = 10;

// let add = fn(x, y) {
//   x + y;
// };

// let result = add(five, ten);
// !-/*5;
// 5 < 10 > 5;

// if (5 < 10) {
// 	return true;
// } else {
// 	return false;
// }

// 10 == 10;
// 10 != 9;
// `

// 	tests := []struct {
// 		expectedType    token.TokenType
// 		expectedLiteral string
// 	}{
// 		{token.LET, "let"},
// 		{token.IDENTIFIER, "five"},
// 		{token.EQUAL, "="},
// 		{token.NUMBER, "5"},
// 		{token.SEMICOLON, ";"},
// 		{token.LET, "let"},
// 		{token.IDENTIFIER, "ten"},
// 		{token.EQUAL, "="},
// 		{token.NUMBER, "10"},
// 		{token.SEMICOLON, ";"},
// 		{token.LET, "let"},
// 		{token.IDENTIFIER, "add"},
// 		{token.EQUAL, "="},
// 		{token.FUNCTION, "fn"},
// 		{token.LEFT_PAREN, "("},
// 		{token.IDENTIFIER, "x"},
// 		{token.COMMA, ","},
// 		{token.IDENTIFIER, "y"},
// 		{token.RIGHT_PAREN, ")"},
// 		{token.LEFT_BRACE, "{"},
// 		{token.IDENTIFIER, "x"},
// 		{token.PLUS, "+"},
// 		{token.IDENTIFIER, "y"},
// 		{token.SEMICOLON, ";"},
// 		{token.RIGHT_BRACE, "}"},
// 		{token.SEMICOLON, ";"},
// 		{token.LET, "let"},
// 		{token.IDENTIFIER, "result"},
// 		{token.EQUAL, "="},
// 		{token.IDENTIFIER, "add"},
// 		{token.LEFT_PAREN, "("},
// 		{token.IDENTIFIER, "five"},
// 		{token.COMMA, ","},
// 		{token.IDENTIFIER, "ten"},
// 		{token.RIGHT_PAREN, ")"},
// 		{token.SEMICOLON, ";"},
// 		{token.BANG, "!"},
// 		{token.MINUS, "-"},
// 		{token.SLASH, "/"},
// 		{token.STAR, "*"},
// 		{token.NUMBER, "5"},
// 		{token.SEMICOLON, ";"},
// 		{token.NUMBER, "5"},
// 		{token.LESS, "<"},
// 		{token.NUMBER, "10"},
// 		{token.GREATER, ">"},
// 		{token.NUMBER, "5"},
// 		{token.SEMICOLON, ";"},
// 		{token.IF, "if"},
// 		{token.LEFT_PAREN, "("},
// 		{token.NUMBER, "5"},
// 		{token.LESS, "<"},
// 		{token.NUMBER, "10"},
// 		{token.RIGHT_PAREN, ")"},
// 		{token.LEFT_BRACE, "{"},
// 		{token.RETURN, "return"},
// 		{token.TRUE, "true"},
// 		{token.SEMICOLON, ";"},
// 		{token.RIGHT_BRACE, "}"},
// 		{token.ELSE, "else"},
// 		{token.LEFT_BRACE, "{"},
// 		{token.RETURN, "return"},
// 		{token.FALSE, "false"},
// 		{token.SEMICOLON, ";"},
// 		{token.RIGHT_BRACE, "}"},
// 		{token.NUMBER, "10"},
// 		{token.EQUAL_EQUAL, "=="},
// 		{token.NUMBER, "10"},
// 		{token.SEMICOLON, ";"},
// 		{token.NUMBER, "10"},
// 		{token.BANG_EQUAL, "!="},
// 		{token.NUMBER, "9"},
// 		{token.SEMICOLON, ";"},
// 		{token.EOF, ""},
// 	}

// 	l := New(input)

// 	for i, tt := range tests {
// 		tok := l.NextToken()

// 		if tok.Type != tt.expectedType {
// 			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
// 				i, tt.expectedType, tok.Type)
// 		}

// 		if tok.Literal != tt.expectedLiteral {
// 			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
// 				i, tt.expectedLiteral, tok.Literal)
// 		}
// 	}
// }
