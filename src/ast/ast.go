package ast

import "meet-lang/src/token"

const (
	NAME              = "NAME"
	INTEGER           = "INTEGER"
	STRING            = "STRING"
	NUMBER            = "NUMBER"
	LIST              = "LIST"
	FUCK_LIST         = "FUCK_LIST"
	EMPTY             = "EMPTY"
	EXP               = "EXP"
	BOOL              = "BOOL"
	MINUS_ONE         = "MINUS_ONE"
	PLUS_ONE          = "PLUS_ONE"
	DEFINE_FUN        = "DEFINE_FUN"
	CALL_FUN          = "CALL_FUN"
	PRINT_SPLACE      = "PRINT_SPLACE"
	PRINT_LINE_SPLACE = "PRINT_LINE_SPLACE"
)

type Program struct {
	Statements []interface{}
}

type FuckStatement struct {
	Name  string
	Type  string
	Value interface{}
}

// ------------------------------------------

type PrintLineStatement struct {
	Type  string
	Value interface{}
}

// ------------------------------------------

type PrintStatement struct {
	Type  string
	Value interface{}
}

// ------------------------------------------

type ListStatement struct {
	Type string
	Size int
	List []interface{}
}

// ------------------------------------------

type BinaryExpressionStatement struct {
	Left     token.Token
	Operator token.Token
	Right    token.Token
}

// ------------------------------------------

type ForEachStatement struct {
	Type string
	Name string
}

// ------------------------------------------

type SetStatement struct {
	Type  string
	Name  string
	Value interface{}
}

// ------------------------------------------

type IfStatement struct {
	Condition []interface{}
	Establish []interface{}
	Contrary  []interface{}
}

// ------------------------------------------

type MinusOnePlusOneStatement struct {
	Type string
	Name string
}

// ------------------------------------------

type WhileStatement struct {
	Condition []interface{}
	Establish []interface{}
}

// ------------------------------------------

type ForStatement struct {
	Establish []interface{}
}

// ------------------------------------------

type BreakStatement struct {
}

// ------------------------------------------

type FunStatement struct {
	Name      string
	Type      string
	Establish []interface{}
}

// ------------------------------------------

type ImportStatement struct {
	Path      string
	Establish []interface{}
}
