package token

type Token struct {
	Type  string
	Value string
}

const (
	NAME   = "NAME"
	DIGIT  = "DIGIT"
	STRING = "STRING"
	LIST   = "LIST"
	BOOL   = "BOOL"

	LPAREN   = "LPAREN"   // (
	RPAREN   = "RPAREN"   // )
	LBRACE   = "LBRACE"   // {
	RBRACE   = "RBRACE"   // }
	LBRACKET = "LBRACKET" // [
	RBRACKET = "RBRACKET" // ]

	PLUS            = "PLUS"            // +
	PLUS_ASSIGN     = "PLUS_ASSIGN"     // +=
	MINUS           = "MINUS"           // -
	MINUS_ASSIGN    = "MINUS_ASSIGN"    // -=
	ASTERISK        = "ASTERISK"        // *
	ASTERISK_ASSIGN = "ASTERISK_ASSIGN" // *=
	DIV             = "DIV"             // /
	DIV_ASSIGN      = "DIV_ASSIGN"      // /=
	EQ              = "EQ"              // ==
	NOT_EQ          = "NOT_EQ"          // !=
	BANG            = "BANG"            // !
	MODULAR         = "MODULAR"         // %
	PLUS_ONE        = "PLUS_ONE"        // ++
	MINUS_ONE       = "MINUS_ONE"       // --

	LT        = "LT"        // <
	LT_ASSIGN = "LT_ASSIGN" // >=
	RT        = "RT"        // >
	RT_ASSIGN = "RT_ASSIGN" // <=

	ASSIGN    = "ASSIGN"    // =
	SEMICOLON = "SEMICOLON" // ;
	COMMA     = "COMMA"     // ,

	POINTER          = "POINTER"          // ->
	FUNCTION_POINTER = "FUNCTION_POINTER" // =>

	EOF = "EOF" // EOF
)
