package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

type Parser struct {
	lexer        *lexer.Lexer
	errors       []string
	currentToken token.Token
	peekToken    token.Token
}

func New(lexer *lexer.Lexer) *Parser {
	parser := &Parser{
		lexer:  lexer,
		errors: []string{},
	}
	parser.nextToken()
	parser.nextToken()
	return parser
}

func (parser *Parser) nextToken() {
	parser.currentToken = parser.peekToken
	parser.peekToken = parser.lexer.NextToken()
}

func (parser *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for parser.currentToken.Type != token.EOF {
		stmt := parser.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		parser.nextToken()
	}

	return program
}

func (parser *Parser) parseStatement() ast.Statement {
	switch parser.currentToken.Type {
	case token.LET:
		return parser.parseLetStatement()
	case token.RETURN:
		return parser.parseReturnStatement()
	default:
		return nil
	}
}

func (parser *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: parser.currentToken}

	if !parser.expectPeek(token.IDENTIFIER) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: parser.currentToken,
		Value: parser.currentToken.Literal}

	if !parser.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: we're skipping the expression until we encounter a semicolon
	for !parser.currentTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return stmt
}

func (parser *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: parser.currentToken}

	parser.nextToken()

	// TODO: we're skipping the expression until we encounter a semicolon
	for !parser.currentTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return stmt
}

func (parser *Parser) currentTokenIs(t token.TokenType) bool {
	return parser.currentToken.Type == t
}

func (parser *Parser) peekTokenIs(t token.TokenType) bool {
	return parser.peekToken.Type == t
}

func (parser *Parser) expectPeek(t token.TokenType) bool {
	if parser.peekTokenIs(t) {
		parser.nextToken()
		return true
	} else {
		parser.peekError(t)
		return false
	}
}

func (parser *Parser) Errors() []string {
	return parser.errors
}

func (parser *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, parser.peekToken.Type)
	parser.errors = append(parser.errors, msg)
}
