package environment

import "bytes"

type Package struct {
	Path      string
	Establish []interface{}
}

func (p *Package) Type() ObjectType {
	return PACKAGE_OBJ
}

func (p *Package) Inspect() string {
	var out bytes.Buffer

	for _, v := range p.Establish {
		out.WriteString(v.(string))
	}

	return out.String()
}

func (p *Package) Body() []interface{} {
	return p.Establish
}
