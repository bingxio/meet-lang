package environment

import (
	"bytes"
	"meet-lang/src/ast"
)

type Fun struct {
	Param     ast.Param
	Establish []interface{}
}

func (f *Fun) Type() ObjectType {
	return FUN_OBJ
}

func (f *Fun) Inspect() string {
	var out bytes.Buffer

	for _, v := range f.Establish {
		out.WriteString(v.(string))
	}

	return out.String()
}

func (f *Fun) Params() ast.Param {
	return f.Param
}

func (f *Fun) Body() []interface{} {
	return f.Establish
}
