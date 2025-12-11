package ast

import (
	"bytes"
	"strings"
	"victoria/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// LetStatement
type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// ReturnStatement
type ReturnStatement struct {
	Token       token.Token // the 'return' token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

// IncludeStatement
type IncludeStatement struct {
	Token   token.Token // the token.INCLUDE token
	Modules []string
}

func (is *IncludeStatement) statementNode()       {}
func (is *IncludeStatement) TokenLiteral() string { return is.Token.Literal }
func (is *IncludeStatement) String() string {
	var out bytes.Buffer
	out.WriteString(is.TokenLiteral() + " ")
	if len(is.Modules) > 1 {
		out.WriteString("(")
		out.WriteString(strings.Join(is.Modules, ", "))
		out.WriteString(")")
	} else if len(is.Modules) == 1 {
		out.WriteString("\"" + is.Modules[0] + "\"")
	}
	out.WriteString(";")
	return out.String()
}

// TryStatement
type TryStatement struct {
	Token      token.Token // token.TRY
	Block      *BlockStatement
	CatchVar   *Identifier
	CatchBlock *BlockStatement
}

func (ts *TryStatement) statementNode()       {}
func (ts *TryStatement) expressionNode()      {}
func (ts *TryStatement) TokenLiteral() string { return ts.Token.Literal }
func (ts *TryStatement) String() string {
	var out bytes.Buffer
	out.WriteString("try ")
	out.WriteString(ts.Block.String())
	if ts.CatchBlock != nil {
		out.WriteString(" catch ")
		if ts.CatchVar != nil {
			out.WriteString("(")
			out.WriteString(ts.CatchVar.String())
			out.WriteString(") ")
		}
		out.WriteString(ts.CatchBlock.String())
	}
	return out.String()
}

// ExpressionStatement
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// BlockStatement
type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// Identifier
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// IntegerLiteral
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// StringLiteral
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

// FloatLiteral
type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }

// Boolean
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

// PrefixExpression
type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. ! or -
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

// InfixExpression
type InfixExpression struct {
	Token    token.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

// IfExpression
type IfExpression struct {
	Token       token.Token // The 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement // Can be BlockStatement or IfExpression (for else if) - simplified to BlockStatement for now, else if handled by parser nesting or list
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

// FunctionLiteral
type FunctionLiteral struct {
	Token      token.Token // The 'define' token
	Name       string      // Optional name, for methods or named functions
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	out.WriteString(fl.TokenLiteral())
	if fl.Name != "" {
		out.WriteString(" " + fl.Name)
	}
	out.WriteString("(")
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())
	return out.String()
}

// CallExpression
type CallExpression struct {
	Token     token.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

// ArrayLiteral
type ArrayLiteral struct {
	Token    token.Token // the '[' token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer
	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

// IndexExpression
type IndexExpression struct {
	Token token.Token // The [ token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")
	return out.String()
}

// HashLiteral
type HashLiteral struct {
	Token token.Token // the '{' token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode()      {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) String() string {
	var out bytes.Buffer
	pairs := []string{}
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

// StructLiteral (Definition)
type StructLiteral struct {
	Token  token.Token // 'struct'
	Name   *Identifier
	Fields []*Identifier
}

func (sl *StructLiteral) statementNode()       {}
func (sl *StructLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StructLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("struct ")
	out.WriteString(sl.Name.String())
	out.WriteString(" { ")
	for _, f := range sl.Fields {
		out.WriteString(f.String() + " ")
	}
	out.WriteString("}")
	return out.String()
}

// PostfixExpression
type PostfixExpression struct {
	Token    token.Token // The operator token, e.g. ++
	Operator string
	Left     Expression
}

func (pe *PostfixExpression) expressionNode()      {}
func (pe *PostfixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PostfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Left.String())
	out.WriteString(pe.Operator)
	out.WriteString(")")
	return out.String()
}

// StructInstantiation
type StructInstantiation struct {
	Token  token.Token // The struct name identifier
	Name   *Identifier
	Fields map[string]Expression
}

func (si *StructInstantiation) expressionNode()      {}
func (si *StructInstantiation) TokenLiteral() string { return si.Token.Literal }
func (si *StructInstantiation) String() string {
	var out bytes.Buffer
	out.WriteString(si.Name.String())
	out.WriteString(" { ")
	pairs := []string{}
	for key, value := range si.Fields {
		pairs = append(pairs, key+": "+value.String())
	}
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString(" }")
	return out.String()
}

// WhileExpression
type WhileExpression struct {
	Token     token.Token
	Condition Expression
	Body      *BlockStatement
}

func (we *WhileExpression) expressionNode()      {}
func (we *WhileExpression) TokenLiteral() string { return we.Token.Literal }
func (we *WhileExpression) String() string {
	var out bytes.Buffer
	out.WriteString("while ")
	out.WriteString(we.Condition.String())
	out.WriteString(" ")
	out.WriteString(we.Body.String())
	return out.String()
}

// ForExpression (Iterating over list/map)
type ForExpression struct {
	Token    token.Token
	Item     *Identifier
	Iterable Expression
	Body     *BlockStatement
}

func (fe *ForExpression) expressionNode()      {}
func (fe *ForExpression) TokenLiteral() string { return fe.Token.Literal }
func (fe *ForExpression) String() string {
	var out bytes.Buffer
	out.WriteString("for ")
	out.WriteString(fe.Item.String())
	out.WriteString(" in ")
	out.WriteString(fe.Iterable.String())
	out.WriteString(" ")
	out.WriteString(fe.Body.String())
	return out.String()
}

// CForExpression
type CForExpression struct {
	Token     token.Token
	Init      Statement
	Condition Expression
	Update    Statement
	Body      *BlockStatement
}

func (cfe *CForExpression) expressionNode()      {}
func (cfe *CForExpression) TokenLiteral() string { return cfe.Token.Literal }
func (cfe *CForExpression) String() string {
	var out bytes.Buffer
	out.WriteString("for (")
	out.WriteString(cfe.Init.String())
	out.WriteString("; ")
	out.WriteString(cfe.Condition.String())
	out.WriteString("; ")
	out.WriteString(cfe.Update.String())
	out.WriteString(") ")
	out.WriteString(cfe.Body.String())
	return out.String()
}

// MethodDefinition
type MethodDefinition struct {
	Token      token.Token // 'define'
	StructName *Identifier
	MethodName *Identifier
	Parameters []*Identifier
	Body       *BlockStatement
}

func (md *MethodDefinition) statementNode()       {}
func (md *MethodDefinition) TokenLiteral() string { return md.Token.Literal }
func (md *MethodDefinition) String() string {
	var out bytes.Buffer
	out.WriteString("define ")
	out.WriteString(md.StructName.String())
	out.WriteString(".")
	out.WriteString(md.MethodName.String())
	out.WriteString("(")
	// params
	out.WriteString(") ")
	out.WriteString(md.Body.String())
	return out.String()
}

// BreakStatement
type BreakStatement struct {
	Token token.Token
}

func (bs *BreakStatement) statementNode()       {}
func (bs *BreakStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BreakStatement) String() string       { return "break" }

// ContinueStatement
type ContinueStatement struct {
	Token token.Token
}

func (cs *ContinueStatement) statementNode()       {}
func (cs *ContinueStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *ContinueStatement) String() string       { return "continue" }

// SwitchExpression
type SwitchExpression struct {
	Token   token.Token
	Value   Expression
	Cases   []*CaseExpression
	Default *BlockStatement
}

func (se *SwitchExpression) expressionNode()      {}
func (se *SwitchExpression) TokenLiteral() string { return se.Token.Literal }
func (se *SwitchExpression) String() string {
	var out bytes.Buffer
	out.WriteString("switch ")
	out.WriteString(se.Value.String())
	out.WriteString(" { ")
	for _, c := range se.Cases {
		out.WriteString(c.String())
	}
	if se.Default != nil {
		out.WriteString("default ")
		out.WriteString(se.Default.String())
	}
	out.WriteString("}")
	return out.String()
}

// CaseExpression
type CaseExpression struct {
	Token token.Token
	Value Expression
	Body  *BlockStatement
}

func (ce *CaseExpression) expressionNode()      {}
func (ce *CaseExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CaseExpression) String() string {
	var out bytes.Buffer
	out.WriteString("case ")
	out.WriteString(ce.Value.String())
	out.WriteString(" ")
	out.WriteString(ce.Body.String())
	return out.String()
}

// TernaryExpression
type TernaryExpression struct {
	Token       token.Token // The ? token
	Condition   Expression
	Consequence Expression
	Alternative Expression
}

func (te *TernaryExpression) expressionNode()      {}
func (te *TernaryExpression) TokenLiteral() string { return te.Token.Literal }
func (te *TernaryExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(te.Condition.String())
	out.WriteString(" ? ")
	out.WriteString(te.Consequence.String())
	out.WriteString(" : ")
	out.WriteString(te.Alternative.String())
	out.WriteString(")")
	return out.String()
}

// RangeExpression
type RangeExpression struct {
	Token token.Token // The .. token
	Start Expression
	End   Expression
}

func (re *RangeExpression) expressionNode()      {}
func (re *RangeExpression) TokenLiteral() string { return re.Token.Literal }
func (re *RangeExpression) String() string {
	var out bytes.Buffer
	out.WriteString(re.Start.String())
	out.WriteString("..")
	out.WriteString(re.End.String())
	return out.String()
}

// ForInIndexExpression - for i, v in arr { }
type ForInIndexExpression struct {
	Token    token.Token
	Index    *Identifier
	Value    *Identifier
	Iterable Expression
	Body     *BlockStatement
}

func (fe *ForInIndexExpression) expressionNode()      {}
func (fe *ForInIndexExpression) TokenLiteral() string { return fe.Token.Literal }
func (fe *ForInIndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("for ")
	out.WriteString(fe.Index.String())
	out.WriteString(", ")
	out.WriteString(fe.Value.String())
	out.WriteString(" in ")
	out.WriteString(fe.Iterable.String())
	out.WriteString(" ")
	out.WriteString(fe.Body.String())
	return out.String()
}
