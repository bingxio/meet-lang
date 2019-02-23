package ast

type ParamItem struct {
	Type  string
	Value interface{}
}

type Param struct {
	Count     int
	ParamItem []ParamItem
}
