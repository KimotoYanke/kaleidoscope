package parser

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"unicode"

	"github.com/KimotoYanke/kaleidoscope/ast"
)

// Token is enum of token
type Token int

const (
	tokEOF Token = iota
	tokDef
	tokExtern
	tokIdentifier
	tokNumber
)

// Parser is a struct to parse kaleidoscope
type Parser struct {
	reader          io.Reader
	numVal          float64
	identifierStr   string
	currentToken    Token
	binopPrecedence map[Token]int
}

func (p *Parser) GetToken() (Token, error) {
	reader := bufio.NewReader(p.reader)
	lastChar := []rune(" ")[0]
	var err error
	var n int

	for unicode.IsSpace(lastChar) {
		lastChar, _, err = reader.ReadRune()
		if err != nil {
			return tokEOF, err
		}
	}

	if unicode.IsLetter(lastChar) {
		identifierStr := []rune{lastChar}

		lastChar, _, err = reader.ReadRune()
		if err == io.EOF {
			return tokEOF, nil
		}
		if err != nil {
			return tokEOF, err
		}
		for unicode.IsLetter(lastChar) || unicode.IsNumber(lastChar) {
			lastChar, n, err = reader.ReadRune()
			if n == 0 {
				return tokEOF, nil
			}
			if err != nil {
				return tokEOF, err
			}

			identifierStr = append(identifierStr, lastChar)

			if reflect.DeepEqual(identifierStr, []rune("def")) {
				return tokDef, nil
			}
			if reflect.DeepEqual(identifierStr, []rune("extern")) {
				return tokExtern, nil
			}
		}
		p.identifierStr = string(identifierStr)
		return tokIdentifier, nil
	}
	if unicode.IsDigit(lastChar) || lastChar == '.' {
		numStr := []rune{}
		for {
			numStr = append(numStr, lastChar)
			lastChar, n, err = reader.ReadRune()
			if err == io.EOF {
				return tokEOF, nil
			}
			if err != nil {
				return tokEOF, err
			}

			if !(unicode.IsDigit(lastChar) || lastChar == rune('.')) {
				break
			}
		}
		numVal, _ := strconv.ParseFloat(string(numStr), 64)
		p.numVal = numVal
		return tokNumber, nil
	}
	if lastChar == '#' {
		reader.ReadLine()
		return p.GetToken()
	}

	if err == io.EOF {
		return tokEOF, nil
	}

	thisChar := lastChar
	lastChar, _, _ = reader.ReadRune()
	return Token(int(thisChar)), nil
}

func (p *Parser) GetNextToken() (Token, error) {
	token, err := p.GetToken()
	if err != nil {
		return -1, nil
	}
	p.currentToken = token
	return p.currentToken, err
}

func (p *Parser) LogError(msg string) ast.ExprAST {
	fmt.Printf(msg)
	return nil
}

func (p *Parser) ParseNumberExpr() ast.ExprAST {
	result := ast.NewNumberExprAST(p.numVal)
	p.GetNextToken()
	return result
}

func (p *Parser) ParseParenExpr() ast.ExprAST {
	p.GetNextToken()
	v := p.ParseExpression()
	if v == nil {
		return nil
	}

	if p.currentToken == ')' {
		return nil
	}

	p.GetNextToken()
	return v
}

func (p *Parser) ParseIdentifierExpr() ast.ExprAST {
	idName := p.identifierStr
	p.GetNextToken()
	if p.currentToken != '(' {
		return ast.NewVariableExprAST(idName).ExprAST
	}
	p.GetNextToken()
	args := []ast.ExprAST{}
	if p.currentToken != ')' {
		for {
			arg := p.ParseExpression()
			if arg != nil {
				args = append(args, arg)
			} else {
				return nil
			}

			if p.currentToken == ')' {
				break
			}

			if p.currentToken != ',' {
				return p.LogError("Expected ')' or ',' in argument list")
			}
			p.GetNextToken()
		}
	}
	p.GetNextToken()

	return ast.NewCallExprAST(idName, args).ExprAST
}

func (p *Parser) ParsePrimary() ast.ExprAST {
	switch p.currentToken {
	case tokIdentifier:
		return p.ParseIdentifierExpr()
	case tokNumber:
		return p.ParseNumberExpr()
	case '(':
		return p.ParseParenExpr()
	default:
		return p.LogError("unknown token when expecting an expression")
	}
}

func isASCII(s Token) bool {
	if s > unicode.MaxASCII {
		return false
	}
	if s < 32 {
		return false
	}
	return true
}

func (p *Parser) GetTokenPrecedence() int {
	if !isASCII(p.currentToken) {
		return -1
	}

	tokenPrecedence := p.binopPrecedence[p.currentToken]
	if tokenPrecedence <= 0 {
		return -1
	}
	return tokenPrecedence
}

func (p *Parser) ParseExpression() ast.ExprAST {
	lhs := p.ParsePrimary()
	if lhs == nil {
		return nil
	}

	return p.ParseBinOpRHS(0, lhs)
}

func (p *Parser) ParseBinOpRHS(exprPrec int, lhs ast.ExprAST) ast.ExprAST {

	for {
		tokenPrec := p.GetTokenPrecedence()
		if tokenPrec < exprPrec {
			return lhs
		}

		binop := p.currentToken
		p.GetNextToken()

		rhs := p.ParsePrimary()
		if rhs == nil {
			return nil
		}

		nextPrec := p.GetTokenPrecedence()
		if tokenPrec < nextPrec {
			movedRhs := rhs
			rhs = p.ParseBinOpRHS(tokenPrec+1, &movedRhs)
			if rhs == nil {
				return nil
			}
		}

		movedLhs := lhs
		movedRhs := rhs

		lhs = ast.NewBinaryExprAST(rune(binop), &movedLhs, &movedRhs).ExprAST
	}
}
