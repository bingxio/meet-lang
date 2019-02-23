package ast

type IfStatement struct {
	Condition []interface{}
	Establish []interface{}
	Contrary  []interface{}
}
