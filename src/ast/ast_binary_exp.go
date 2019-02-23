package ast

import "meet-lang/src/token"

type BinaryExpressionStatement struct {
	Left     token.Token
	Operator token.Token
	Right    token.Token
}
