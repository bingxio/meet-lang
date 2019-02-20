package main

import (
	"flag"
	"fmt"
	"meet-lang/src/config"
	"meet-lang/src/environment"
	"meet-lang/src/interpreter"
	"meet-lang/src/lexer"
	"meet-lang/src/parser"
	"meet-lang/src/repl"
	"meet-lang/src/util"
	"os"
	"os/user"
)

var env = environment.NewEnvironment()

func main() {
	user, err := user.Current()

	if err != nil {
		panic(err)
	}

	shellIsShowTokens := flag.Bool("token", false, "show program tokens")
	shellIsShowAst := flag.Bool("ast", false, "show program ast")
	shellIsShowEnv := flag.Bool("env", false, "show program env")
	shellIsShowAll := flag.Bool("all", false, "show tokens and ast and env")
	shellIsShowMore := flag.Bool("more", false, "show more about meet programming language")

	flag.Parse()

	shellWithFile := flag.Arg(0)

	if *shellIsShowMore {
		fmt.Println(config.MORE)
		return
	}

	if shellWithFile == "" {
		fmt.Printf("Hello %s! Meet Programming Language REPL %s - Turaiiao 2019 - Email: 1171840237@qq.com\n", user.Username, config.VERSION)

		repl.Start(os.Stdin, os.Stdout, env)

		return
	}

	if !util.Suffix(shellWithFile) {
		panic("文件仅限 .meet 后缀")
	} else {
		code, err := util.Read(shellWithFile)

		if err != nil {
			panic("文件打开失败：" + err.Error() + ", path = " + shellWithFile)
		}

		eval(string(code), *shellIsShowTokens, *shellIsShowAst, *shellIsShowEnv, *shellIsShowAll)
	}
}

func eval(code string, show_token, show_ast, show_env bool, show_all bool) {
	l := lexer.New(code)

	if show_token {
		for k, v := range l {
			fmt.Println(k, v)
		}
		splitScreen()
	}

	ast := parser.New(l).ParseProgram()

	if show_ast {
		for k, v := range ast.Statements {
			fmt.Println(k, v)
		}
		splitScreen()
	}

	interpreter.Eval(ast, env)

	if show_env {
		splitScreen()

		for k, v := range env.All() {
			fmt.Printf("k = %s, v = %d \n", k, v)
		}
	}

	if show_all {
		splitScreen()

		for k, v := range l {
			fmt.Println(k, v)
		}

		splitScreen()

		for k, v := range ast.Statements {
			fmt.Println(k, v)
		}

		splitScreen()

		for k, v := range env.All() {
			fmt.Printf("k = %s, v = %s \n", k, v)
		}
	}
}

func splitScreen() {
	for i := 0; i < 30; i++ {
		fmt.Print("-")
	}
	fmt.Println()
}
