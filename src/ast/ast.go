package ast

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
