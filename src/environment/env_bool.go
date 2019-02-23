package environment

import "strconv"

type Bool struct {
	State bool
}

func (b *Bool) Type() ObjectType {
	return BOOL_OBJ
}

func (b *Bool) Inspect() string {
	return strconv.FormatBool(b.State)
}
