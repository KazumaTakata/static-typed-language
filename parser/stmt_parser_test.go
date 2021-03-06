package parser

import (
	"fmt"
	"github.com/KazumaTakata/regex_virtualmachine"
	"github.com/KazumaTakata/static-typed-language/lexer"
	"testing"
)

func TestStmtParser(t *testing.T) {

	fmt.Printf("\n------------------------------------------------------\n")

	regex_string := lexer.Get_Regex_String()

	regex := regex.NewRegexWithParser(regex_string)

	input := "var x int = 3+13.0 \n\n\n\n x \n\n"

	tokens := lexer.GetTokens(regex, input)

	parser_input := Parser_Input{Tokens: tokens, Pos: 0}

	fmt.Printf("%+v\n\n", tokens)

	stmt := Parse_Stmts(&parser_input)

	fmt.Printf("%+v", stmt)
}

func TestForstmtParser(t *testing.T) {

	fmt.Printf("\n------------------------------------------------------\n")

	regex_string := lexer.Get_Regex_String()

	regex := regex.NewRegexWithParser(regex_string)

	input := "for x == 3+13.0 {x} \n\n"

	tokens := lexer.GetTokens(regex, input)

	parser_input := Parser_Input{Tokens: tokens, Pos: 0}

	fmt.Printf("%+v\n\n", tokens)

	stmt := Parse_Stmt(&parser_input)

	fmt.Printf("%+v", stmt.For)
}

func TestDefStmtParser(t *testing.T) {

	fmt.Printf("\n------------------------------------------------------\n")

	regex_string := lexer.Get_Regex_String()

	regex := regex.NewRegexWithParser(regex_string)

	input := "def add(a int, b int) int { a } \n\n"

	tokens := lexer.GetTokens(regex, input)

	parser_input := Parser_Input{Tokens: tokens, Pos: 0}

	fmt.Printf("%+v\n\n", tokens)

	stmt := Parse_Stmt(&parser_input)

	fmt.Printf("%+v", stmt.Def)
}

func TestCallParser(t *testing.T) {

	fmt.Printf("\n------------------------------------------------------\n")

	regex_string := lexer.Get_Regex_String()

	regex := regex.NewRegexWithParser(regex_string)

	input := "add(a, b) \n\n"

	tokens := lexer.GetTokens(regex, input)

	parser_input := Parser_Input{Tokens: tokens, Pos: 0}

	fmt.Printf("%+v\n\n", tokens)

	stmt := Parse_Stmt(&parser_input)

	fmt.Printf("%+v", stmt.Expr)
}

func TestAssignParser(t *testing.T) {

	fmt.Printf("\n------------------------------------------------------\n")

	regex_string := lexer.Get_Regex_String()

	regex := regex.NewRegexWithParser(regex_string)

	input := "a = 3 + e \n\n"

	tokens := lexer.GetTokens(regex, input)

	parser_input := Parser_Input{Tokens: tokens, Pos: 0}

	fmt.Printf("%+v\n\n", tokens)

	stmt := Parse_Stmt(&parser_input)

	fmt.Printf("%+v", stmt.Assign)
}

func TestInitParser(t *testing.T) {

	fmt.Printf("\n------------------------------------------------------\n")

	regex_string := lexer.Get_Regex_String()

	regex := regex.NewRegexWithParser(regex_string)

	input := "var a []int =[]int{3, 3}  \n\n"

	tokens := lexer.GetTokens(regex, input)

	parser_input := Parser_Input{Tokens: tokens, Pos: 0}

	fmt.Printf("%+v\n\n", tokens)

	stmt := Parse_Stmt(&parser_input)

	fmt.Printf("%+v", stmt.Decl)
}

func TestNestedArray(t *testing.T) {

	fmt.Printf("\n------------------------------------------------------\n")

	regex_string := lexer.Get_Regex_String()

	regex := regex.NewRegexWithParser(regex_string)

	input := "var a [][]int = [][]int{[]int{3}}  \n\n"

	tokens := lexer.GetTokens(regex, input)

	parser_input := Parser_Input{Tokens: tokens, Pos: 0}

	fmt.Printf("%+v\n\n", tokens)

	stmt := Parse_Stmt(&parser_input)

	fmt.Printf("%+v", stmt.Decl)
}

func TestFor(t *testing.T) {

	fmt.Printf("\n------------------------------------------------------\n")

	regex_string := lexer.Get_Regex_String()

	regex := regex.NewRegexWithParser(regex_string)

	input := "for var i int = 0; i < 10; i = i + 1 {} \n\n"

	tokens := lexer.GetTokens(regex, input)

	parser_input := Parser_Input{Tokens: tokens, Pos: 0}

	fmt.Printf("%+v\n\n", tokens)

	stmt := Parse_Stmt(&parser_input)

	fmt.Printf("%+v", stmt.For)
}
