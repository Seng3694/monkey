package lexer

import "monkey/token"

type Lexer struct {
	input        string
	position     int
	readPosition int
	char         byte
}

func New(input string) *Lexer {
	lexer := &Lexer{input: input}
	lexer.readChar()
	return lexer
}

func (lexer *Lexer) NextToken() token.Token {
	var tok token.Token

	lexer.skipWhitespace()

	switch lexer.char {
	case '=':
		if lexer.peekChar() == '=' {
			lexer.readChar()
			tok = token.Token{Type: token.EQ, Literal: "=="}
		} else {
			tok = newToken(token.ASSIGN, lexer.char)
		}
	case '+':
		tok = newToken(token.PLUS, lexer.char)
	case '-':
		tok = newToken(token.MINUS, lexer.char)
	case '!':
		if lexer.peekChar() == '=' {
			lexer.readChar()
			tok = token.Token{Type: token.NEQ, Literal: "!="}
		} else {
			tok = newToken(token.BANG, lexer.char)
		}
	case '/':
		tok = newToken(token.SLASH, lexer.char)
	case '*':
		tok = newToken(token.ASTERISK, lexer.char)
	case '<':
		if lexer.peekChar() == '=' {
			lexer.readChar()
			tok = token.Token{Type: token.LEQ, Literal: "<="}
		} else {
			tok = newToken(token.LT, lexer.char)
		}
	case '>':
		if lexer.peekChar() == '=' {
			lexer.readChar()
			tok = token.Token{Type: token.GEQ, Literal: ">="}
		} else {
			tok = newToken(token.GT, lexer.char)
		}
	case ',':
		tok = newToken(token.COMMA, lexer.char)
	case ';':
		tok = newToken(token.SEMICOLON, lexer.char)
	case '(':
		tok = newToken(token.LPAREN, lexer.char)
	case ')':
		tok = newToken(token.RPAREN, lexer.char)
	case '{':
		tok = newToken(token.LBRACE, lexer.char)
	case '}':
		tok = newToken(token.RBRACE, lexer.char)
	case '[':
		tok = newToken(token.LBRACKET, lexer.char)
	case ']':
		tok = newToken(token.RBRACKET, lexer.char)
	case '"':
		tok.Type = token.STRING
		tok.Literal = lexer.readString()
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(lexer.char) {
			tok.Literal = lexer.readIdentifier()
			tok.Type = token.LookupIdentifier(tok.Literal)
			return tok
		} else if isDigit(lexer.char) {
			tok.Literal = lexer.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, lexer.char)
		}
	}
	lexer.readChar()
	return tok
}

func newToken(tokenType token.TokenType, char byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(char)}
}

func isWhitespace(char byte) bool {
	return char == ' ' || char == '\t' || char == '\n' || char == '\r'
}

func isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

func isLetter(char byte) bool {
	return char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' || char == '_'
}

func (lexer *Lexer) skipWhitespace() {
	for isWhitespace(lexer.char) {
		lexer.readChar()
	}
}

func (lexer *Lexer) peekChar() byte {
	if lexer.readPosition >= len(lexer.input) {
		return 0
	} else {
		return lexer.input[lexer.readPosition]
	}
}

func (lexer *Lexer) readChar() {
	if lexer.readPosition >= len(lexer.input) {
		lexer.char = 0
	} else {
		lexer.char = lexer.input[lexer.readPosition]
	}
	lexer.position = lexer.readPosition
	lexer.readPosition += 1
}

func (lexer *Lexer) readIdentifier() string {
	position := lexer.position
	if isLetter(lexer.char) {
		for isLetter(lexer.char) || isDigit(lexer.char) {
			lexer.readChar()
		}
	}

	return lexer.input[position:lexer.position]
}

func (lexer *Lexer) readNumber() string {
	position := lexer.position
	for isDigit(lexer.char) {
		lexer.readChar()
	}
	return lexer.input[position:lexer.position]
}

func (lexer *Lexer) readString() string {
	position := lexer.position + 1
	for {
		lexer.readChar()
		if lexer.char == '"' || lexer.char == 0 {
			break
		}
	}
	return lexer.input[position:lexer.position]
}
