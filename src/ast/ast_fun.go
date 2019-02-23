package ast

type FunStatement struct {
	Name      string
	Type      string
	Param     Param
	Establish []interface{}
}
