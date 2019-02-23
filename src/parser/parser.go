package parser

import (
	"fmt"
	"meet-lang/src/ast"
	"meet-lang/src/lexer"
	"meet-lang/src/token"
	"meet-lang/src/util"
	"regexp"
	"strconv"
)

type Parser struct {
	tokens  []token.Token
	token   token.Token
	current int
	length  int
}

func New(tokens []token.Token) *Parser {
	p := Parser{
		tokens:  tokens,
		current: 0,
	}
	return &p
}

func (p *Parser) ParseProgram() *ast.Program {
	p.length = len(p.tokens)

	program := &ast.Program{}

	for p.current < p.length && p.currentToken().Type != token.EOF {
		p.token = p.tokens[p.current]
		p.parseWithToken(&program.Statements)
	}

	return program
}

func (p *Parser) parseWithToken(statements *[]interface{}) {
	switch p.token.Value {
	case "import":
		*statements = append(*statements, *p.parseImportStatement())
	case "fuck":
		*statements = append(*statements, *p.parseFuckStatement())
	case "print":
		*statements = append(*statements, *p.parsePrintStatement())
	case "printLine":
		*statements = append(*statements, *p.parsePrintLineStatement())
	case "forEach":
		*statements = append(*statements, *p.parseForEachStatement())
	case "set":
		*statements = append(*statements, *p.parseSetStatement())
	case "if":
		*statements = append(*statements, *p.parseIfStatement())
	case "while":
		*statements = append(*statements, *p.parseWhileStatement())
	case "break":
		*statements = append(*statements, *p.parseBreakStatement())
	case "for":
		*statements = append(*statements, *p.parseForStatement())
	case "fun":
		*statements = append(*statements, *p.parseFunStatement())
	default:
		if p.token.Type == token.NAME && p.tokens[p.current+1].Type == token.MINUS_ONE ||
			p.tokens[p.current+1].Type == token.PLUS_ONE {
			*statements = append(*statements, *p.parseMinusOnePlusOneStatement())
		} else if p.token.Type == token.NAME && p.tokens[p.current+1].Type == token.POINTER {
			*statements = append(*statements, *p.parseReFuckStatement())
		} else {
			panic("未知的类型：" + p.token.Type + ", " + p.token.Value)
		}
	}
}

func (p *Parser) parseExpressionStatement() (string, interface{}) {
	p.refreshCurrentToken()

	t := p.currentToken()

	switch t.Type {
	case token.DIGIT:
		v, _ := strconv.Atoi(t.Value)
		return ast.INTEGER, v
	case token.NAME:
		return ast.NAME, t.Value
	case token.STRING:
		return ast.STRING, t.Value
	case token.LBRACKET:
		return ast.LIST, *p.parseListStatement()
	case token.LPAREN:
		return ast.EXP, *p.parseBinaryExpressionStatement()
	case token.LIST:
		return ast.FUCK_LIST, t.Value
	}

	panic("未知的表达式类型：" + t.Type + ", " + t.Value)
}

func (p *Parser) parseBlockStatement(ast []interface{}) *[]interface{} {
	for p.currentToken().Type != token.RBRACE {
		p.parseWithToken(&ast)
		p.refreshCurrentToken()
	}

	p.current++
	p.refreshCurrentToken() // refresh token '}' to next token.

	return &ast
}

func (p *Parser) parseParamsStatement() *ast.Param {
	params := &ast.Param{}

	for p.currentToken().Type != token.RPAREN {
		if p.currentToken().Type == token.COMMA {
			p.current++
			continue
		}

		params.ParamItem = append(params.ParamItem, ast.ParamItem{
			Type:  p.currentToken().Type,
			Value: p.currentToken().Value,
		})

		p.current++
		p.refreshCurrentToken()
	}

	p.current++
	p.refreshCurrentToken() // skip ')'

	params.Count = len(params.ParamItem)

	return params
}

// import 'sample_package.meet' as sample;
func (p *Parser) parseImportStatement() *ast.ImportStatement {
	importStmt := &ast.ImportStatement{}

	p.current++

	if p.currentToken().Type != token.STRING {
		panic("缺少包名")
	}

	importStmt.Path = p.currentToken().Value

	p.current++
	p.isSemicolon()
	p.current++
	p.refreshCurrentToken()

	if !util.Suffix(importStmt.Path) {
		panic("文件仅限 .meet 后缀")
	}

	code, err := util.Read(importStmt.Path)

	if err != nil {
		panic("文件打开失败：" + err.Error())
	}

	tokens := lexer.New(string(code))
	ast := New(tokens).ParseProgram()

	importStmt.Establish = ast.Statements

	return importStmt
}

// fuck a -> 20;
// fuck a -> 'meet programming language !';
// fuck a -> [2 4 6 8 10];
// fuck a -> (1 + 2);
// fuck a -> (list[0] + 3);
// fuck a -> list[0];
func (p *Parser) parseFuckStatement() *ast.FuckStatement {
	fuckStmt := &ast.FuckStatement{}

	name := p.nextToken()

	p.isLetter(name)
	p.current++
	p.isPointer()
	p.current++

	types, value := p.parseExpressionStatement()

	fuckStmt.Name = name.Value
	fuckStmt.Value = value
	fuckStmt.Type = types

	if types != ast.EXP {
		p.current++
	}

	p.isSemicolon()
	p.current++
	p.refreshCurrentToken()

	return fuckStmt
}

// print -> a;
// print -> 'meet programming language !';
// print;
// print 2;
// print index;
// print -> list[0];
// print -> (2 + 2);
// print -> (list[0] * 2);
// print -> (2 == 2);
func (p *Parser) parsePrintStatement() *ast.PrintStatement {
	printStmt := &ast.PrintStatement{}

	p.current++

	if p.currentToken().Type == token.SEMICOLON {
		printStmt.Type = ast.EMPTY

		p.current++
		p.refreshCurrentToken()

		return printStmt
	} else if p.currentToken().Type == token.DIGIT {
		printStmt.Type = ast.NUMBER
		printStmt.Value = p.currentToken().Value

		p.current++
		p.isSemicolon()
		p.current++
		p.refreshCurrentToken()

		return printStmt
	} else if p.currentToken().Type == token.NAME {
		printStmt.Type = ast.PRINT_SPLACE
		printStmt.Value = p.currentToken().Value

		p.current++
		p.isSemicolon()
		p.current++
		p.refreshCurrentToken()

		return printStmt
	}

	p.isPointer()
	p.current++

	if p.currentToken().Type == token.LIST {
		printStmt.Type = ast.LIST
		printStmt.Value = p.currentToken().Value

		p.current++
		p.isSemicolon()
		p.current++
		p.refreshCurrentToken()

		return printStmt
	}

	types, value := p.parseExpressionStatement()

	printStmt.Type = types
	printStmt.Value = value

	if types != ast.EXP {
		p.current++
	}

	p.isSemicolon()
	p.current++
	p.refreshCurrentToken()

	return printStmt
}

// printLine -> a;
// printLine -> 'meet programming language !';
// printLine;
// printLine 2;
// printLine index;
// printLine -> list[0];
// printLine -> (2 + 2);
// printLine -> (list[0] * 2);
// printLine -> (2 == 2);
func (p *Parser) parsePrintLineStatement() *ast.PrintLineStatement {
	printLineStmt := &ast.PrintLineStatement{}

	p.current++

	if p.currentToken().Type == token.SEMICOLON {
		printLineStmt.Type = ast.EMPTY

		p.current++
		p.refreshCurrentToken()

		return printLineStmt
	} else if p.currentToken().Type == token.DIGIT {
		printLineStmt.Type = ast.NUMBER
		printLineStmt.Value = p.currentToken().Value

		p.current++
		p.isSemicolon()
		p.current++
		p.refreshCurrentToken()

		return printLineStmt
	} else if p.currentToken().Type == token.NAME {
		printLineStmt.Type = ast.PRINT_LINE_SPLACE
		printLineStmt.Value = p.currentToken().Value

		p.current++
		p.isSemicolon()
		p.current++
		p.refreshCurrentToken()

		return printLineStmt
	}

	p.isPointer()
	p.current++

	if p.currentToken().Type == token.LIST {
		printLineStmt.Type = ast.LIST
		printLineStmt.Value = p.currentToken().Value

		p.current++
		p.isSemicolon()
		p.current++
		p.refreshCurrentToken()

		return printLineStmt
	}

	types, value := p.parseExpressionStatement()

	printLineStmt.Type = types
	printLineStmt.Value = value

	if types != ast.EXP {
		p.current++
	}

	p.isSemicolon()
	p.current++
	p.refreshCurrentToken()

	return printLineStmt
}

// forEach -> list;
func (p *Parser) parseForEachStatement() *ast.ForEachStatement {
	forEachStmt := &ast.ForEachStatement{}

	p.current++
	p.isPointer()
	p.current++
	p.isLetter(p.currentToken())

	_, value := p.parseExpressionStatement()

	forEachStmt.Name = value.(string)
	forEachStmt.Type = ast.LIST

	p.current++
	p.isSemicolon()
	p.current++
	p.refreshCurrentToken()

	return forEachStmt
}

// set list[0] -> 20;
func (p *Parser) parseSetStatement() *ast.SetStatement {
	setStmt := &ast.SetStatement{}

	p.current++

	if p.currentToken().Type != token.LIST {
		panic("找不到列表名：" + p.currentToken().Value)
	}

	setStmt.Name = p.currentToken().Value

	p.current++
	p.isPointer()
	p.current++

	types, value := p.parseExpressionStatement()

	setStmt.Type = types
	setStmt.Value = value

	if types != ast.EXP {
		p.current++
	}

	p.isSemicolon()
	p.current++
	p.refreshCurrentToken()

	return setStmt
}

// a ++; a --;
func (p *Parser) parseMinusOnePlusOneStatement() *ast.MinusOnePlusOneStatement {
	minusOnePlusOneStmt := &ast.MinusOnePlusOneStatement{}
	minusOnePlusOneStmt.Name = p.currentToken().Value

	p.current++
	p.refreshCurrentToken()

	if p.currentToken().Type == token.MINUS_ONE {
		minusOnePlusOneStmt.Type = ast.MINUS_ONE
	} else if p.currentToken().Type == token.PLUS_ONE {
		minusOnePlusOneStmt.Type = ast.PLUS_ONE
	} else {
		panic("不是位加或位减操作：" + p.currentToken().Value)
	}

	p.current++
	p.isSemicolon()
	p.current++
	p.refreshCurrentToken()

	return minusOnePlusOneStmt
}

// if (a += 1) > b {
// 	printLine -> a;
// } else {
// 	println -> b;
// }
func (p *Parser) parseIfStatement() *ast.IfStatement {
	ifStmt := &ast.IfStatement{}

	p.current++

	for p.currentToken().Type != token.LBRACE {
		if p.currentToken().Type == token.LPAREN {
			ifStmt.Condition = append(ifStmt.Condition, *p.parseBinaryExpressionStatement())
			p.refreshCurrentToken()
		}
		ifStmt.Condition = append(ifStmt.Condition, p.currentToken())
		p.current++
		p.refreshCurrentToken()
	}

	p.isLBrace()
	p.current++
	p.refreshCurrentToken()

	ifStmt.Establish = *p.parseBlockStatement(ifStmt.Establish)

	if p.currentToken().Value == "else" {
		p.current++ // skip 'else'
		p.refreshCurrentToken()

		if p.currentToken().Type != token.LBRACE {
			panic("缺少大括号：" + p.currentToken().Value)
		}

		p.current++ // skip '{'
		p.refreshCurrentToken()

		ifStmt.Contrary = *p.parseBlockStatement(ifStmt.Contrary)
	}

	return ifStmt
}

// while (a += 1) > b {
//     printLine -> a;
// }
func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	whileStmt := &ast.WhileStatement{}

	p.current++

	for p.currentToken().Type != token.LBRACE {
		if p.currentToken().Type == token.LPAREN {
			whileStmt.Condition = append(whileStmt.Condition, *p.parseBinaryExpressionStatement())
			p.refreshCurrentToken()
		}
		whileStmt.Condition = append(whileStmt.Condition, p.currentToken())
		p.current++
		p.refreshCurrentToken()
	}

	p.isLBrace()
	p.current++
	p.refreshCurrentToken()

	whileStmt.Establish = *p.parseBlockStatement(whileStmt.Establish)

	return whileStmt
}

// break;
func (p *Parser) parseBreakStatement() *ast.BreakStatement {
	breakStmt := &ast.BreakStatement{}

	p.current++
	p.isSemicolon()
	p.current++
	p.refreshCurrentToken()

	return breakStmt
}

// for {
//     break;
// }
func (p *Parser) parseForStatement() *ast.ForStatement {
	forStmt := &ast.ForStatement{}

	p.current++
	p.isLBrace()
	p.current++
	p.refreshCurrentToken()

	forStmt.Establish = *p.parseBlockStatement(forStmt.Establish)

	return forStmt
}

// [2 3 4 5 6 7 8] -> ast.ListStatement { Value: []interface }
// ['meet' 'programming' 'language'] ast.INTEGER / ast.STRING
func (p *Parser) parseListStatement() *ast.ListStatement {
	listStmt := &ast.ListStatement{}
	listStmt.Size = 0

	p.current++ // skip '['

	for !p.isToken("]") {
		t := p.currentToken()

		// skip ','
		if t.Type == token.COMMA {
			p.current++
			continue
		}

		if listStmt.Type == "" {
			listStmt.Type = t.Type
		} else if t.Type != listStmt.Type {
			panic("列表类型必须一致：" + t.Type + ", " + t.Value)
		}

		listStmt.List = append(listStmt.List, t.Value)
		listStmt.Size++

		p.current++
	}

	if len(listStmt.List) == 0 {
		listStmt.Type = ast.INTEGER
	}

	// 先判断的列表值类型、现在设置 ast 值类型
	if listStmt.Type == token.DIGIT {
		listStmt.Type = ast.INTEGER
	} else if listStmt.Type == token.STRING {
		listStmt.Type = ast.STRING
	}

	return listStmt
}

// (2 / 4) -> ast.BinaryExpressionStatement {DIGIT 2} {DIV /} {DIGIT 4}
func (p *Parser) parseBinaryExpressionStatement() *ast.BinaryExpressionStatement {
	binaryExpressionStatement := &ast.BinaryExpressionStatement{}
	binaryExpressionList := make([]interface{}, 0)

	p.current++ // skip '('

	for p.currentToken().Type != token.RPAREN {
		binaryExpressionList = append(binaryExpressionList, p.currentToken())
		p.current++
	}

	if len(binaryExpressionList) > 3 {
		panic("操作数大于 3 个：" + strconv.Itoa(len(binaryExpressionList)))
	}

	binaryExpressionStatement.Left = binaryExpressionList[0].(token.Token)
	binaryExpressionStatement.Operator = binaryExpressionList[1].(token.Token)
	binaryExpressionStatement.Right = binaryExpressionList[2].(token.Token)

	p.current++ // skip ')'

	return binaryExpressionStatement
}

// fun a => (a, b) {
// 	printLine -> a;
// }
//
// fun -> a (12, 20);
func (p *Parser) parseFunStatement() *ast.FunStatement {
	funStmt := &ast.FunStatement{}

	p.current++
	p.refreshCurrentToken()

	if p.currentToken().Type != token.POINTER {
		p.isLetter(p.currentToken())

		funStmt.Name = p.currentToken().Value

		p.current++

		if p.currentToken().Type == token.FUNCTION_POINTER {
			p.current++

			if p.currentToken().Type == token.LPAREN {
				p.current++ // skip '('

				funStmt.Param = *p.parseParamsStatement()
			}

			p.isLBrace()
			p.current++
			p.refreshCurrentToken()

			funStmt.Type = ast.DEFINE_FUN
			funStmt.Establish = *p.parseBlockStatement(funStmt.Establish)
		} else {
			panic("缺少函数指针：" + p.currentToken().Value)
		}

		return funStmt
	} else {
		p.isPointer()
		p.current++
		p.isLetter(p.currentToken())

		funStmt.Name = p.currentToken().Value

		p.current++

		if p.currentToken().Type == token.LPAREN {
			p.current++ // skip '('

			funStmt.Param = *p.parseParamsStatement()
			funStmt.Type = ast.CALL_FUN

			p.isSemicolon()
			p.current++
			p.refreshCurrentToken()

			return funStmt
		}

		funStmt.Type = ast.CALL_FUN

		p.refreshCurrentToken()
		p.isSemicolon()
		p.current++
		p.refreshCurrentToken()
	}

	return funStmt
}

// a -> 20;
func (p *Parser) parseReFuckStatement() *ast.ReFuckStatement {
	reFuckStmt := &ast.ReFuckStatement{}

	p.isLetter(p.currentToken())

	reFuckStmt.Name = p.currentToken().Value

	p.current++
	p.isPointer()
	p.current++

	types, value := p.parseExpressionStatement()

	reFuckStmt.Type = types
	reFuckStmt.Value = value

	p.current++
	p.isSemicolon()
	p.current++
	p.refreshCurrentToken()

	return reFuckStmt
}

// -------------------------------------------

func (p Parser) showCurrentToken() {
	fmt.Println(p.token)
}

func (p Parser) showCurrentTokens() {
	fmt.Println(p.tokens[p.current])
}

func (p Parser) currentToken() token.Token {
	return p.tokens[p.current]
}

func (p *Parser) nextToken() token.Token {
	p.current++

	if p.currentToken().Type == token.EOF {
		panic("语法错误：" + p.currentToken().Value)
	}

	return p.tokens[p.current]
}

func (p Parser) isToken(value string) bool {
	return p.tokens[p.current].Value == value
}

func (p *Parser) refreshCurrentToken() {
	p.token = p.tokens[p.current]
}

func (p Parser) isLetter(t token.Token) bool {
	r, _ := regexp.Compile("[a-z|A-Z]")

	if !r.MatchString(t.Value) {
		panic("变量名只能是英文字符序列：" + p.currentToken().Value)
	}

	return true
}

func (p Parser) isPointer() bool {
	if p.currentToken().Type != token.POINTER {
		panic("缺少变量指针：" + p.currentToken().Value)
	}

	return true
}

func (p Parser) isSemicolon() bool {
	if p.currentToken().Type != token.SEMICOLON {
		panic("缺少分号：" + p.currentToken().Value)
	}

	return true
}

func (p Parser) isLBrace() bool {
	if p.currentToken().Type != token.LBRACE {
		panic("缺少左大括号：" + p.currentToken().Value)
	}

	return true
}

func (p Parser) isRBrace() bool {
	if p.currentToken().Type != token.RBRACE {
		panic("缺少右大括号：" + p.currentToken().Value)
	}

	return true
}
