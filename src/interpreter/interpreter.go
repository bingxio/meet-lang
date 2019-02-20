package interpreter

import (
	"fmt"
	"meet-lang/src/ast"
	"meet-lang/src/environment"
	"meet-lang/src/token"
	"regexp"
	"strconv"
)

type Interpreter struct {
	ast     ast.Program
	node    interface{}
	env     environment.Environment
	current int
	length  int

	breakForStatement bool
}

func Eval(ast *ast.Program, env *environment.Environment) {
	i := &Interpreter{
		ast:     *ast,
		env:     *env,
		current: 0,

		breakForStatement: true,
	}
	i.eval()
}

func (i *Interpreter) eval() {
	i.length = len(i.ast.Statements)

	for i.current < i.length {
		i.node = i.ast.Statements[i.current]
		i.evalForNode(i.node)
	}
}

func (i *Interpreter) evalForNode(node interface{}) {
	switch node.(type) {
	case ast.ImportStatement:
		i.evalImportStatementNode()
	case ast.FuckStatement:
		i.evalFuckStatementNode()
	case ast.PrintStatement:
		i.evalPrintStatementNode()
	case ast.PrintLineStatement:
		i.evalPrintLineStatementNode()
	case ast.ForEachStatement:
		i.evalForEachStatementNode()
	case ast.SetStatement:
		i.evalSetStatementNode()
	case ast.IfStatement:
		i.evalIfStatementNode()
	case ast.MinusOnePlusOneStatement:
		i.evalMinusOnePlusOneStatementNode()
	case ast.WhileStatement:
		i.evalWhileStatementNode()
	case ast.BreakStatement:
		i.evalBreakStatement()
	case ast.ForStatement:
		i.evalForStatement()
	case ast.FunStatement:
		i.evalFunStatement()
	default:
		panic("解释失败，未知类型")
	}
}

func (i *Interpreter) evalImportStatementNode() {
	importStmt := i.node.(ast.ImportStatement)
	tempCurrent := i.current

	for _, v := range importStmt.Establish {
		i.node = v
		i.evalForNode(i.node)
	}

	i.current = tempCurrent
	i.current++
}

func (i *Interpreter) evalFuckStatementNode() {
	fuckStmt := i.node.(ast.FuckStatement)

	if fuckStmt.Type == ast.INTEGER {
		i.env.Set(fuckStmt.Name, &environment.Integer{
			Value: fuckStmt.Value.(int),
		})
	} else if fuckStmt.Type == ast.STRING {
		i.env.Set(fuckStmt.Name, &environment.String{
			Value: fuckStmt.Value.(string),
		})
	} else if fuckStmt.Type == ast.LIST {
		v := fuckStmt.Value.(ast.ListStatement)

		// 列表默认值类型转换成一致的
		for t := 0; t < v.Size; t++ {
			if v.Type == environment.INTEGER_OBJ {
				v.List[t] = i.toInt(v.List[t].(string))
			}
		}

		for f := v.Size; f < 1000; f++ {
			if v.Type == ast.STRING {
				v.List = append(v.List, "null")
			} else {
				v.List = append(v.List, 0)
			}
		}

		i.env.Set(fuckStmt.Name, &environment.List{
			Types:  v.Type,
			Size:   v.Size,
			Values: v.List,
		})
	} else if fuckStmt.Type == ast.EXP {
		t, v := i.evalBinaryExpressionNode(fuckStmt.Value.(ast.BinaryExpressionStatement))

		if t == ast.STRING {
			i.env.Set(fuckStmt.Name, &environment.String{
				Value: v.(string),
			})
		} else if t == ast.INTEGER {
			i.env.Set(fuckStmt.Name, &environment.Integer{
				Value: v.(int),
			})
		} else if t == ast.BOOL {
			i.env.Set(fuckStmt.Name, &environment.Bool{
				State: v.(bool),
			})
		}
	} else if fuckStmt.Type == ast.FUCK_LIST {
		_, _, t, _, v := i.evalListExpression(fuckStmt.Value.(string))

		if t == environment.INTEGER_OBJ {
			i.env.Set(fuckStmt.Name, &environment.Integer{
				Value: v.(int),
			})
		} else if t == environment.STRING_OBJ {
			i.env.Set(fuckStmt.Name, &environment.String{
				Value: v.(string),
			})
		}
	}

	i.current++
}

func (i *Interpreter) evalPrintStatementNode() {
	printStmt := i.node.(ast.PrintStatement)

	if printStmt.Type == ast.NAME {
		if v := i.envGetVariable(printStmt.Value.(string)); v.Type() == environment.INTEGER_OBJ ||
			v.Type() == environment.STRING_OBJ || v.Type() == environment.BOOL_OBJ {
			print(v.Inspect())
		} else if v.Type() == environment.LIST_OBJ {
			list := v.(*environment.List)

			print(list.Items())
		}
	} else if printStmt.Type == ast.STRING {
		print(printStmt.Value.(string))
	} else if printStmt.Type == ast.EMPTY {
		print(" ")
	} else if printStmt.Type == ast.NUMBER {
		v, _ := strconv.Atoi(printStmt.Value.(string))

		for f := 0; f < v; f++ {
			print(" ")
		}
	} else if printStmt.Type == ast.INTEGER {
		print(printStmt.Value)
	} else if printStmt.Type == ast.LIST {
		_, _, _, _, list_value := i.evalListExpression(printStmt.Value.(string))

		print(list_value)
	} else if printStmt.Type == ast.EXP {
		_, v := i.evalBinaryExpressionNode(printStmt.Value.(ast.BinaryExpressionStatement))

		print(v)
	} else if printStmt.Type == ast.PRINT_SPLACE {
		v := i.envGetVariable(printStmt.Value.(string))

		for i := 0; i < v.(*environment.Integer).Value; i++ {
			print(" ")
		}
	}

	i.current++
}

func (i *Interpreter) evalPrintLineStatementNode() {
	printLineStmt := i.node.(ast.PrintLineStatement)

	if printLineStmt.Type == ast.NAME {
		if v := i.envGetVariable(printLineStmt.Value.(string)); v.Type() == environment.INTEGER_OBJ ||
			v.Type() == environment.STRING_OBJ || v.Type() == environment.BOOL_OBJ {
			printLine(v.Inspect())
		} else if v.Type() == environment.LIST_OBJ {
			list := v.(*environment.List)

			printLine(list.Items())
		}
	} else if printLineStmt.Type == ast.STRING {
		printLine(printLineStmt.Value.(string))
	} else if printLineStmt.Type == ast.EMPTY {
		printLine("")
	} else if printLineStmt.Type == ast.NUMBER {
		v, _ := strconv.Atoi(printLineStmt.Value.(string))

		for f := 0; f < v; f++ {
			printLine("")
		}
	} else if printLineStmt.Type == ast.INTEGER {
		printLine(printLineStmt.Value)
	} else if printLineStmt.Type == ast.LIST {
		_, _, _, _, list_value := i.evalListExpression(printLineStmt.Value.(string))

		printLine(list_value)
	} else if printLineStmt.Type == ast.EXP {
		_, v := i.evalBinaryExpressionNode(printLineStmt.Value.(ast.BinaryExpressionStatement))

		printLine(v)
	} else if printLineStmt.Type == ast.PRINT_LINE_SPLACE {
		v := i.envGetVariable(printLineStmt.Value.(string))

		for i := 0; i < v.(*environment.Integer).Value; i++ {
			printLine("")
		}
	}

	i.current++
}

func (i *Interpreter) evalForEachStatementNode() {
	forEachStmt := i.node.(ast.ForEachStatement)

	v := i.envGetVariable(forEachStmt.Name)
	l := v.(*environment.List)

	for f := 0; f < l.Size; f++ {
		print(l.Items()[f])
		print(" ")
	}

	printLine("")

	i.current++
}

func (i *Interpreter) evalSetStatementNode() {
	setStmt := i.node.(ast.SetStatement)

	listName, listIdx := i.findListMore(setStmt.Name)

	v := i.envGetVariable(listName)
	l := v.(*environment.List)

	if listIdx > 999 {
		panic("列表最大下标为 999：" + string(listIdx))
	}

	if setStmt.Type == ast.EXP {
		t, v := i.evalBinaryExpressionNode(setStmt.Value.(ast.BinaryExpressionStatement))

		setStmt.Type = t
		setStmt.Value = v
	}

	if setStmt.Type == ast.NAME {
		v := i.envGetVariable(setStmt.Value.(string))

		if v.Type() == environment.INTEGER_OBJ {
			t := v.(*environment.Integer)

			setStmt.Type = ast.INTEGER
			setStmt.Value = t.Value
		} else if v.Type() == environment.STRING_OBJ {
			t := v.(*environment.String)

			setStmt.Type = ast.STRING
			setStmt.Value = t.Value
		}
	}

	if l.Types != setStmt.Type {
		panic("列表值类型需要一致，原列表：" + l.Types + "，新值：" + setStmt.Type)
	}

	l.Items()[listIdx] = setStmt.Value
	l.Size++

	i.env.Set(listName, &environment.List{
		Types:  l.Types,
		Size:   l.Size,
		Values: l.Items(),
	})

	i.current++
}

func (i *Interpreter) evalMinusOnePlusOneStatementNode() {
	minusOnePlusOneStmt := i.node.(ast.MinusOnePlusOneStatement)

	if v := i.envGetVariable(minusOnePlusOneStmt.Name); v.Type() == environment.INTEGER_OBJ {
		i := v.(*environment.Integer)

		if minusOnePlusOneStmt.Type == ast.PLUS_ONE {
			i.Value++
		} else {
			i.Value--
		}
	} else {
		panic("位加位减操作只能对整型运算：" + v.Type())
	}

	i.current++
}

func (i *Interpreter) evalIfStatementNode() {
	ifStmt := i.node.(ast.IfStatement)

	condition := i.evalConditionStatement(ifStmt.Condition)
	_, v := i.evalBinaryExpressionNode(*condition)

	tempCurrent := i.current

	if v.(bool) {
		for _, v := range ifStmt.Establish {
			i.node = v
			i.evalForNode(i.node)
		}

		i.current = tempCurrent
	} else {
		for _, v := range ifStmt.Contrary {
			i.node = v
			i.evalForNode(i.node)
		}

		i.current = tempCurrent
	}

	i.current++
}

func (i *Interpreter) evalWhileStatementNode() {
	whileStmt := i.node.(ast.WhileStatement)

	condition := i.evalConditionStatement(whileStmt.Condition)
	_, v := i.evalBinaryExpressionNode(*condition)

	tempCurrent := i.current

	for v.(bool) {
		for _, n := range whileStmt.Establish {
			i.node = n
			i.evalForNode(i.node)
		}

		condition = i.evalConditionStatement(whileStmt.Condition)
		_, v = i.evalBinaryExpressionNode(*condition)
	}

	i.current = tempCurrent

	i.current++
}

func (i *Interpreter) evalBreakStatement() {
	i.breakForStatement = false

	i.current++
}

func (i *Interpreter) evalForStatement() {
	forStmt := i.node.(ast.ForStatement)

	tempCurrent := i.current

	for i.breakForStatement {
		for _, n := range forStmt.Establish {
			i.node = n
			i.evalForNode(i.node)
		}
	}

	i.breakForStatement = true

	i.current = tempCurrent
	i.current++
}

func (i *Interpreter) evalFunStatement() {
	funStmt := i.node.(ast.FunStatement)

	if funStmt.Type == ast.DEFINE_FUN {
		i.env.Set(funStmt.Name, &environment.Fun{
			Establish: funStmt.Establish,
		})
	} else {
		tempCurrent := i.current

		v := i.envGetVariable(funStmt.Name)
		f := v.(*environment.Fun)

		for _, v := range f.Body() {
			i.node = v
			i.evalForNode(i.node)
		}

		i.current = tempCurrent
	}

	i.current++
}

func (i Interpreter) evalConditionStatement(conditionArr []interface{}) *ast.BinaryExpressionStatement {
	condition := ast.BinaryExpressionStatement{}

	if len(conditionArr) > 3 {
		panic("操作数不能大于 3 个：" + string(len(conditionArr)))
	}

	switch conditionArr[0].(type) {
	case ast.BinaryExpressionStatement:
		exp := conditionArr[0].(ast.BinaryExpressionStatement)
		_, value := i.evalBinaryExpressionNode(exp)

		if exp.Operator.Type == token.PLUS_ASSIGN || exp.Operator.Type == token.MINUS_ASSIGN ||
			exp.Operator.Type == token.ASTERISK_ASSIGN || exp.Operator.Type == token.DIV_ASSIGN {
			i.env.Set(exp.Left.Value, &environment.Integer{
				Value: value.(int),
			})
		}

		condition.Left = token.Token{
			Type:  token.DIGIT,
			Value: strconv.Itoa(value.(int)),
		}

		condition.Operator = token.Token{
			Type:  conditionArr[1].(token.Token).Type,
			Value: conditionArr[1].(token.Token).Value,
		}

		condition.Right = token.Token{
			Type:  conditionArr[2].(token.Token).Type,
			Value: conditionArr[2].(token.Token).Value,
		}

		return &condition
	case token.Token:
		condition.Left = token.Token{
			Type:  conditionArr[0].(token.Token).Type,
			Value: conditionArr[0].(token.Token).Value,
		}

		condition.Operator = token.Token{
			Type:  conditionArr[1].(token.Token).Type,
			Value: conditionArr[1].(token.Token).Value,
		}

		condition.Right = token.Token{
			Type:  conditionArr[2].(token.Token).Type,
			Value: conditionArr[2].(token.Token).Value,
		}

		return &condition
	default:
		panic("未知的条件表达式")
	}
}

func (i Interpreter) envGetVariable(name string) environment.Object {
	v, ok := i.env.Get(name)

	if !ok {
		panic("找不到变量：" + name)
	}

	return v
}

// list_name, list_index, list_type, list_size, list_value
func (i Interpreter) evalListExpression(value string) (string, int, string, int, interface{}) {
	listName, listIdx := i.findListMore(value)

	v := i.envGetVariable(listName)
	l := v.(*environment.List)

	if listIdx > (l.Size - 1) {
		panic("越界啦：" + strconv.Itoa(listIdx) + ", 最大下标: " + strconv.Itoa(l.Size-1))
	}

	return listName, listIdx, l.Types, l.Size - 1, l.Items()[listIdx]
}

func (i Interpreter) findListMore(value string) (string, int) {
	listName := ""
	listIndex := ""
	current := 0

	length := len([]rune(value))

	for current < length && string([]rune(value)[current]) != "[" {
		listName += string([]rune(value)[current])
		current++
	}

	current++

	for current < length && string([]rune(value)[current]) != "]" {
		listIndex += string([]rune(value)[current])
		current++
	}

	if listIndex == "" {
		panic("无法获取到列表下标：" + value)
	}

	if i.isLetter(listIndex) {
		v := i.envGetVariable(listIndex)
		listIndex := v.(*environment.Integer).Value

		return listName, listIndex
	}

	listIdx, _ := strconv.Atoi(listIndex)

	return listName, listIdx
}

// ast.INTEGER / ast.STRING -> int / string / bool
func (i Interpreter) evalBinaryExpressionNode(binaryExpressionStatement ast.BinaryExpressionStatement) (string, interface{}) {
	var binaryExpressionStatementLeft interface{}
	var binaryExpressionStatementLeftT string
	var binaryExpressionStatementOperator string
	var binaryExpressionStatementRight interface{}
	var binaryExpressionStatementRightT string

	if binaryExpressionStatement.Left.Type == token.DIGIT {
		binaryExpressionStatementLeft = i.toInt(binaryExpressionStatement.Left.Value)
		binaryExpressionStatementLeftT = ast.INTEGER
	} else if binaryExpressionStatement.Left.Type == token.STRING {
		binaryExpressionStatementLeft = binaryExpressionStatement.Left.Value
		binaryExpressionStatementLeftT = ast.STRING
	} else if binaryExpressionStatement.Left.Type == token.LIST {
		_, _, list_type, _, list_value := i.evalListExpression(binaryExpressionStatement.Left.Value)

		if list_type == environment.INTEGER_OBJ {
			binaryExpressionStatementLeft = list_value.(int)
			binaryExpressionStatementLeftT = ast.INTEGER
		} else if list_type == environment.STRING_OBJ {
			binaryExpressionStatementLeft = list_value.(string)
			binaryExpressionStatementLeftT = ast.STRING
		}
	} else if binaryExpressionStatement.Left.Type == token.NAME {
		if v := i.envGetVariable(binaryExpressionStatement.Left.Value); v.Type() == environment.STRING_OBJ {
			binaryExpressionStatementLeft = v.Inspect()
			binaryExpressionStatementLeftT = ast.STRING
		} else if v.Type() == environment.INTEGER_OBJ {
			binaryExpressionStatementLeft = i.toInt(v.Inspect())
			binaryExpressionStatementLeftT = ast.INTEGER
		}
	}

	if binaryExpressionStatement.Right.Type == token.DIGIT {
		binaryExpressionStatementRight = i.toInt(binaryExpressionStatement.Right.Value)
		binaryExpressionStatementRightT = ast.INTEGER
	} else if binaryExpressionStatement.Right.Type == token.STRING {
		binaryExpressionStatementRight = binaryExpressionStatement.Right.Value
		binaryExpressionStatementRightT = ast.STRING
	} else if binaryExpressionStatement.Right.Type == token.LIST {
		_, _, list_type, _, list_value := i.evalListExpression(binaryExpressionStatement.Right.Value)

		if list_type == environment.INTEGER_OBJ {
			binaryExpressionStatementRight = list_value.(int)
			binaryExpressionStatementRightT = ast.INTEGER
		} else if list_type == environment.STRING_OBJ {
			binaryExpressionStatementRight = list_value.(string)
			binaryExpressionStatementRightT = ast.STRING
		}
	} else if binaryExpressionStatement.Right.Type == token.NAME {
		if v := i.envGetVariable(binaryExpressionStatement.Right.Value); v.Type() == environment.STRING_OBJ {
			binaryExpressionStatementRight = v.Inspect()
			binaryExpressionStatementRightT = ast.STRING
		} else if v.Type() == environment.INTEGER_OBJ {
			binaryExpressionStatementRight = i.toInt(v.Inspect())
			binaryExpressionStatementRightT = ast.INTEGER
		}
	}

	binaryExpressionStatementOperator = binaryExpressionStatement.Operator.Type

	switch binaryExpressionStatementOperator {
	case token.PLUS:
		return i.plus(binaryExpressionStatementLeftT, binaryExpressionStatementRightT, binaryExpressionStatementLeft, binaryExpressionStatementRight)
	case token.MINUS:
		return i.minus(binaryExpressionStatementLeftT, binaryExpressionStatementRightT, binaryExpressionStatementLeft, binaryExpressionStatementRight)
	case token.ASTERISK:
		return i.asterisk(binaryExpressionStatementLeftT, binaryExpressionStatementRightT, binaryExpressionStatementLeft, binaryExpressionStatementRight)
	case token.DIV:
		return i.div(binaryExpressionStatementLeftT, binaryExpressionStatementRightT, binaryExpressionStatementLeft, binaryExpressionStatementRight)
	case token.MODULAR:
		return i.modular(binaryExpressionStatementLeftT, binaryExpressionStatementRightT, binaryExpressionStatementLeft, binaryExpressionStatementRight)
	case token.LT, token.RT, token.LT_ASSIGN, token.RT_ASSIGN, token.EQ, token.NOT_EQ:
		return i.logical(binaryExpressionStatementLeftT, binaryExpressionStatementOperator, binaryExpressionStatementRightT,
			binaryExpressionStatementLeft, binaryExpressionStatementRight)
	case token.PLUS_ASSIGN, token.MINUS_ASSIGN, token.ASTERISK_ASSIGN, token.DIV_ASSIGN:
		return i.opAssign(binaryExpressionStatementLeftT, binaryExpressionStatementOperator, binaryExpressionStatementRightT,
			binaryExpressionStatementLeft, binaryExpressionStatementRight)
	}

	panic("表达式运算出错：left = " + binaryExpressionStatementLeftT +
		", operator = " + binaryExpressionStatementOperator + ", right = " + binaryExpressionStatementRightT)
}

// -------------------------------------------

func print(value interface{}) {
	fmt.Print(value)
}

func printLine(value interface{}) {
	fmt.Println(value)
}

func (i Interpreter) showCurrentNode() {
	fmt.Println(i.node)
}

func (i Interpreter) showCurrentNodes() {
	fmt.Println(i.ast.Statements[i.current])
}

func (i Interpreter) toInt(value string) int {
	v, _ := strconv.Atoi(value)

	return v
}

func (i Interpreter) plus(leftT, rightT string, left, right interface{}) (string, interface{}) {
	if leftT == ast.STRING && rightT == ast.STRING {
		return ast.STRING, left.(string) + right.(string)
	}

	if leftT == ast.INTEGER && rightT == ast.INTEGER {
		return ast.INTEGER, left.(int) + right.(int)
	}

	panic("相加只能是两个字符串或者两个整型：left = " + leftT + " , right = " + rightT)
}

func (i Interpreter) minus(leftT, rightT string, left, right interface{}) (string, interface{}) {
	if leftT == ast.INTEGER && rightT == ast.INTEGER {
		return ast.INTEGER, left.(int) - right.(int)
	}

	panic("相减只能是两个整型：left = " + leftT + ", right = " + rightT)
}

func (i Interpreter) asterisk(leftT, rightT string, left, right interface{}) (string, interface{}) {
	if leftT == ast.INTEGER && rightT == ast.INTEGER {
		return ast.INTEGER, left.(int) * right.(int)
	}

	panic("相乘只能是两个整型：left = " + leftT + ", right = " + rightT)
}

func (i Interpreter) div(leftT, rightT string, left, right interface{}) (string, interface{}) {
	if leftT == ast.INTEGER && rightT == ast.INTEGER {
		return ast.INTEGER, left.(int) / right.(int)
	}

	panic("相除只能是两个整型：left = " + leftT + ", right = " + rightT)
}

func (i Interpreter) modular(leftT, rightT string, left, right interface{}) (string, interface{}) {
	if leftT == ast.INTEGER && rightT == ast.INTEGER {
		return ast.INTEGER, left.(int) % right.(int)
	}

	panic("取模只能是两个整型：left = " + leftT + ", right = " + rightT)
}

func (i Interpreter) logical(leftT, operator, rightT string, left, right interface{}) (string, interface{}) {
	if leftT != ast.INTEGER && rightT != ast.INTEGER {
		panic("逻辑运算类型出错：left = " + leftT + ", right = " + rightT)
	}

	switch operator {
	case token.LT:
		return ast.BOOL, left.(int) > right.(int)
	case token.RT:
		return ast.BOOL, left.(int) < right.(int)
	case token.LT_ASSIGN:
		return ast.BOOL, left.(int) >= right.(int)
	case token.RT_ASSIGN:
		return ast.BOOL, left.(int) <= right.(int)
	case token.EQ:
		return ast.BOOL, left.(int) == right.(int)
	case token.NOT_EQ:
		return ast.BOOL, left.(int) != right.(int)
	default:
		panic("未知的逻辑操作数：" + operator)
	}
}

func (i Interpreter) opAssign(leftT, operator, rightT string, left, right interface{}) (string, interface{}) {
	if leftT != ast.INTEGER && rightT != ast.INTEGER {
		panic("运算类型出错：left = " + leftT + ", right = " + rightT)
	}

	switch operator {
	case token.PLUS_ASSIGN, token.PLUS:
		return ast.INTEGER, left.(int) + right.(int)
	case token.MINUS_ASSIGN, token.MINUS:
		return ast.INTEGER, left.(int) - right.(int)
	case token.ASTERISK_ASSIGN, token.ASTERISK:
		return ast.INTEGER, left.(int) * right.(int)
	case token.DIV_ASSIGN, token.DIV:
		return ast.INTEGER, left.(int) / right.(int)
	default:
		panic("未知的操作数：" + operator)
	}
}

func (i Interpreter) isLetter(v string) bool {
	r, _ := regexp.Compile("[a-z|A-Z]")

	return r.MatchString(v)
}
