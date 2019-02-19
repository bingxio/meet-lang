package environment

import "bytes"

// --------- List ---------

type List struct {
	Types  string
	Size   int
	Values []interface{}
}

func (l *List) Type() ObjectType {
	return LIST_OBJ
}

func (l *List) Inspect() string {
	// return List{Types: l.Types, Size: l.Size, Values: l.Values}
	var out bytes.Buffer

	for _, v := range l.Values {
		out.WriteString(v.(string))
	}

	return out.String()
}

func (l *List) Items() []interface{} {
	return l.Values
}
