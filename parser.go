package parser

import (
	"github.com/KimotoYanke/kaleidoscope/ast"
	parsec "github.com/prataprc/goparsec"
)

func one2one(ns []parsec.ParsecNode) parsec.ParsecNode {
	if ns == nil || len(ns) == 0 {
		return nil
	}
	return ns[0]
}

var openParen = parsec.Token(`\(`, "OPEN_PAREN")
var closeParen = parsec.Token(`\)`, "CLOSE_PAREN")
var addOp = parsec.Token(`+`, "ADD_OP")
var subOp = parsec.Token(`-`, "SUB_OP")
var mulOp = parsec.Token(`*`, "MUL_OP")

var number = parsec.Float()

var sumOp = parsec.OrdChoice(one2one, addOp, subOp)

// Nodifiers

func binNode(ns []parsec.ParsecNode) parsec.ParsecNode {
	if len(ns) == 3 {
		op := ns[0]
		lhs := ns[1]
		rhs := ns[2]
		ast.NewBinaryExprAST(op, lhs, rhs)

	}
}
