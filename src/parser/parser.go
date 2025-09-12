package parser

import (
	"fmt"

	"github.com/looksaw/interpreter/src/ast"
	"github.com/looksaw/interpreter/src/lexer"
	"github.com/looksaw/interpreter/src/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX 
	CALL
)


type Parser struct {
	l *lexer.Lexer
	curToken token.Token
	peekToken token.Token
	errors []string
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFn map[token.TokenType]infixParseFn
}


func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l :l,
		errors: []string{},
	}
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT,p.parseIdentifier)
	p.infixParseFn = make(map[token.TokenType]infixParseFn)
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser)Errors() []string {
	return p.errors
}
func (p *Parser)peekError(t token.TokenType){
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",t,p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser)ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}


func (p *Parser)parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token:p.curToken}
	if !p.expectPeek(token.IDENT){
		return nil
	}
	stmt.Name = &ast.Identifier{
		Token : p.curToken,
		Value: p.curToken.Literal,
	}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	if !p.curTokenIs(token.SEMICOLON){
		p.nextToken()
	}
	return stmt
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser)peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser)expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}else {
		p.peekError(t)
		return false
	}
}


func (p *Parser)parseReturnStatement() *ast.ReturnStatemnet {
	stmt := &ast.ReturnStatemnet{Token: p.curToken}
	p.nextToken()

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser)parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	//stmt.Expression = p.parseExpression()
	stmt.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func(p *Parser)parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		return nil
	}
	leftExp := prefix()
	return leftExp
}

func (p *Parser)parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

//递归下降
type(
	prefixParseFn func() ast.Expression
	infixParseFn func(ast.Expression) ast.Expression
)

//前缀函数注册
func (p *Parser) registerPrefix(tokenType token.TokenType , fn prefixParseFn){
	p.prefixParseFns[tokenType] = fn
}
//中缀函数注册
func (p *Parser)registerInfix(tokenType token.TokenType , fn infixParseFn){
	p.infixParseFn[tokenType] = fn
}