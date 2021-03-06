package eval

import (
	"fmt"
	"github.com/KazumaTakata/regex_virtualmachine"
	"github.com/KazumaTakata/static-typed-language/lexer"
	"github.com/KazumaTakata/static-typed-language/parser"
	basic_type "github.com/KazumaTakata/static-typed-language/type"
	"github.com/KazumaTakata/static-typed-language/type-system"
	"io/ioutil"
)

func Eval_Stmts(stmts []parser.Stmt, symbol_env *parser.Symbol_Env) {

	for _, stmt := range stmts {
		if_return := Eval_Stmt(stmt, symbol_env)

		if if_return {
			break
		}
	}
}

func Eval_Init(init parser.Init, symbol_env *parser.Symbol_Env) parser.Object {

	switch init.Type {
	case parser.ARRAY_INIT:
		{
			arrayobj := parser.ArrayObj{ElementType: init.Array.ElementType}

			for _, init_value := range init.Array.InitValue {
				assign := Eval_Assign(*init_value, symbol_env)
				arrayobj.Value = append(arrayobj.Value, &assign)

			}

			return parser.Object{Type: parser.ArrayType, Array: &arrayobj}
		}
	case parser.MAP_INIT:
		{
		}
	}

	return parser.Object{}
}

func Eval_Assign(assign parser.Assign, symbol_env *parser.Symbol_Env) parser.Object {
	switch assign.Type {
	case parser.INIT_ASSIGN:
		{
			return Eval_Init(*assign.Init, symbol_env)
		}
	case parser.EXPR_ASSIGN:
		{
			return Calc_Arith(assign.Expr.Left, symbol_env)
		}
	}

	return parser.Object{}
}

func Calc_Arith(expr *parser.Arith_expr, symbol_env *parser.Symbol_Env) parser.Object {

	switch expr.Type.DataStructType {
	case basic_type.PRIMITIVE:
		{
			switch expr.Type.Primitive.Type {
			case basic_type.INT:
				{
					result := Arith_Terms_INT(expr.Terms, symbol_env)
					return parser.Object{Type: parser.PrimitiveType, Primitive: &parser.PrimitiveObj{Int: result, Type: basic_type.INT}}

				}
			case basic_type.DOUBLE:
				{
					result := Arith_Terms_DOUBLE(expr.Terms, symbol_env)
					return parser.Object{Type: parser.PrimitiveType, Primitive: &parser.PrimitiveObj{Double: result, Type: basic_type.DOUBLE}}

				}
			case basic_type.STRING:
				{
					result := Arith_Terms_STRING(expr.Terms, symbol_env)
					return parser.Object{Type: parser.PrimitiveType, Primitive: &parser.PrimitiveObj{String: result, Type: basic_type.STRING}}

				}

			case basic_type.BOOL:
				{
					result := Arith_Terms_BOOL(expr.Terms, symbol_env)
					return parser.Object{Type: parser.PrimitiveType, Primitive: &parser.PrimitiveObj{Bool: result, Type: basic_type.BOOL}}

				}

			}
		}

	case basic_type.ARRAY:
		{

			return *resolve_ident(expr.Terms[0].Term.Factors[0].Factor, symbol_env)
		}

	default:
		{
			return *resolve_ident(expr.Terms[0].Term.Factors[0].Factor, symbol_env)
		}
	}

	return parser.Object{}
}

func assign_Table(id string, symbol_env *parser.Symbol_Env, object parser.Object) {
	if _, ok := symbol_env.Table[id]; ok {
		symbol_env.Table[id] = object
	} else {
		assign_Table(id, symbol_env.Parent_Env, object)
	}

}

func assign_to_array_object(object *parser.Object, indexs []int, assign_object parser.Object) *parser.Object {
	if len(indexs) == 1 {
		object.Array.Value[indexs[0]] = &assign_object
		return object
	}
	object.Array.Value[indexs[0]] = assign_to_array_object(object.Array.Value[indexs[0]], indexs[1:], assign_object)

	return object

}

func eval_Assign(assign *parser.Assign_stmt, symbol_env *parser.Symbol_Env) {

	result := Eval_Assign(*assign.Assign, symbol_env)

	if len(assign.Indexs) > 0 {
		//index := Calc_Arith(&stmt.Assign.Indexs[0], symbol_env)
		object := resolve_variable(assign.Id, symbol_env)
		indexs := []int{}
		for _, index := range assign.Indexs {
			index := Calc_Arith(&index, symbol_env)
			indexs = append(indexs, index.Primitive.Int)
		}
		object = *assign_to_array_object(&object, indexs, result)
		symbol_env.Table[assign.Id] = object

	} else {
		assign_Table(assign.Id, symbol_env, result)
	}

}

func Eval_Stmt(stmt parser.Stmt, symbol_env *parser.Symbol_Env) bool {

	switch stmt.Type {

	case parser.DEF_STMT:
		{

		}

	case parser.IMPORT_STMT:
		{
			regex_string := lexer.Get_Regex_String()

			regex := regex.NewRegexWithParser(regex_string)

			module_symbol_env := parser.Symbol_Env{Table: parser.Symbol_Table{}}

			dat, _ := ioutil.ReadFile(stmt.Import.Module_name + ".cat")
			string_input := string(dat)

			tokens := lexer.GetTokens(regex, string_input)
			parser_input := parser.Parser_Input{Tokens: tokens, Pos: 0}
			stmts := parser.Parse_Stmts(&parser_input)

			type_system.Type_Check_Stmts(stmts, &module_symbol_env)

			Eval_Stmts(stmts, &module_symbol_env)

			symbol_env.Table[stmt.Import.Module_name] = parser.Object{Type: parser.EnvType, Env: &module_symbol_env}

		}
	case parser.RETURN_STMT:
		{
			return_value := Calc_Arith(stmt.Return.Cmp_expr.Left, symbol_env)
			symbol_env.Return_Value = &return_value

			return true

		}

	case parser.FOR_STMT:
		{
			switch stmt.For.Type {
			case parser.Cmp:
				{
					for Eval_Cmp_Bool(stmt.For.Cmp_expr, symbol_env) {
						Eval_Stmts(stmt.For.Stmts, stmt.For.Symbol_Env)
					}
				}
			case parser.DeclCmpAssign:
				{

					result := Eval_Assign(*stmt.For.Decl.Assign, stmt.For.Symbol_Env)
					stmt.For.Symbol_Env.Table[stmt.For.Decl.Id] = result

					for Eval_Cmp_Bool(stmt.For.Cmp_expr, stmt.For.Symbol_Env) {
						Eval_Stmts(stmt.For.Stmts, stmt.For.Symbol_Env)
						eval_Assign(&stmt.For.Assign, stmt.For.Symbol_Env)
					}

				}
			}

		}
	case parser.IF_STMT:
		{
			if Eval_Cmp_Bool(stmt.If.Cmp_expr, symbol_env) {
				Eval_Stmts(stmt.If.Stmts, stmt.If.Symbol_Env)
			}
		}

	case parser.ASSIGN_STMT:
		{

			result := Eval_Assign(*stmt.Assign.Assign, symbol_env)

			if len(stmt.Assign.Indexs) > 0 {
				//index := Calc_Arith(&stmt.Assign.Indexs[0], symbol_env)
				object := resolve_variable(stmt.Assign.Id, symbol_env)
				indexs := []int{}
				for _, index := range stmt.Assign.Indexs {
					index := Calc_Arith(&index, symbol_env)
					indexs = append(indexs, index.Primitive.Int)
				}
				object = *assign_to_array_object(&object, indexs, result)
				symbol_env.Table[stmt.Assign.Id] = object

			} else {
				assign_Table(stmt.Assign.Id, symbol_env, result)
			}
		}

	case parser.DECL_STMT:
		{

			result := Eval_Assign(*stmt.Decl.Assign, symbol_env)
			symbol_env.Table[stmt.Decl.Id] = result
			//fmt.Printf("%+v", result)

		}

	case parser.EXPR:
		{

			switch stmt.Expr.Type.DataStructType {
			case basic_type.PRIMITIVE:
				{

					switch stmt.Expr.Type.Primitive.Type {
					case basic_type.INT:
						{
							result := Arith_Terms_INT(stmt.Expr.Terms, symbol_env)
							fmt.Printf("%+v\n", result)

						}
					case basic_type.DOUBLE:
						{
							result := Arith_Terms_DOUBLE(stmt.Expr.Terms, symbol_env)

							fmt.Printf("%+v\n", result)

						}
					case basic_type.STRING:
						{
							result := Arith_Terms_STRING(stmt.Expr.Terms, symbol_env)

							fmt.Printf("%+v\n", result)

						}
					case basic_type.BOOL:
						{
							result := Arith_Terms_BOOL(stmt.Expr.Terms, symbol_env)

							fmt.Printf("%+v\n", result)

						}

					}
				}
			case basic_type.ARRAY:
				{
					factor := stmt.Expr.Terms[0].Term.Factors[0].Factor

					object := resolve_ident(factor, symbol_env)
					parser.PrintObject(*object)

				}
			default:
				{
					factor := stmt.Expr.Terms[0].Term.Factors[0].Factor

					_ = resolve_ident(factor, symbol_env)

				}

			}
		}
	}
	return false
}
