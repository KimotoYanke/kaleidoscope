package ast

type ExprAST interface {
}

type NumberExprAST struct {
	ExprAST
	val float64
}

func NewNumberExprAST(val float64) *NumberExprAST {
	return &NumberExprAST{val: val}
}

type VariableExprAST struct {
	ExprAST
	name string
}

func NewVariableExprAST(name string) *VariableExprAST {
	return &VariableExprAST{name: name}
}

type BinaryExprAST struct {
	ExprAST
	op  rune
	lhs *ExprAST
	rhs *ExprAST
}

func NewBinaryExprAST(op rune, lhs *ExprAST, rhs *ExprAST) *BinaryExprAST {
	return &BinaryExprAST{op: op, lhs: lhs, rhs: rhs}
}

type CallExprAST struct {
	ExprAST
	callee string
	args   []*ExprAST
}

func NewCallExprAST(callee string, args []*ExprAST) *CallExprAST {
	return &CallExprAST{callee: callee, args: args}
}

type PrototypeAST struct {
	ExprAST
	name string
	args []string
}

func NewPrototypeAST(name string, args []string) *PrototypeAST {
	return &PrototypeAST{name: name, args: args}
}

type FunctionAST struct {
	ExprAST
	proto *PrototypeAST
	body  *ExprAST
}

func NewFunctionAST(proto *PrototypeAST, body *ExprAST) *FunctionAST {
	return &FunctionAST{proto: proto, body: body}
}
