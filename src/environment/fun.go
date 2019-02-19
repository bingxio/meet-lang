package environment

import "bytes"

type Fun struct {
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

func (f *Fun) Body() []interface{} {
	return f.Establish
}
