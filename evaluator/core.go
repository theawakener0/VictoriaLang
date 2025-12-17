package evaluator

import (
	"victoria/ast"
	"victoria/lexer"
	"victoria/object"
	"victoria/parser"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// Eval evaluates an AST node and returns the result
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.IncludeStatement:
		return evalIncludeStatement(node, env)

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		// Type check if type annotation is present
		if node.Type != nil {
			if !object.CheckType(val, node.Type) {
				return newErrorWithLocation("type mismatch: cannot assign %s to variable of type %s",
					node.Token.Line, node.Token.Column, node.Token.EndColumn+len(node.Name.Value),
					object.TypeName(val), node.Type.String())
			}
		}
		env.Set(node.Name.Value, val)

	case *ast.ConstStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		// Type check if type annotation is present
		if node.Type != nil {
			if !object.CheckType(val, node.Type) {
				return newErrorWithLocation("type mismatch: cannot assign %s to constant of type %s",
					node.Token.Line, node.Token.Column, node.Token.EndColumn+len(node.Name.Value),
					object.TypeName(val), node.Type.String())
			}
		}
		env.SetConst(node.Name.Value, val)

	case *ast.MakeStatement:
		// #make defines a compile-time constant (like C's #define)
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.SetConst(node.Name.Value, val)

	case *ast.EnumStatement:
		// Create enum type and register all values
		enumObj := &object.Enum{
			Name:   node.Name.Value,
			Values: make(map[string]int64),
		}
		var nextValue int64 = 0
		for _, v := range node.Values {
			var value int64
			if v.Value != nil {
				valObj := Eval(v.Value, env)
				if isError(valObj) {
					return valObj
				}
				if intObj, ok := valObj.(*object.Integer); ok {
					value = intObj.Value
					nextValue = value + 1
				} else {
					return newError("enum value must be an integer")
				}
			} else {
				value = nextValue
				nextValue++
			}
			enumObj.Values[v.Name.Value] = value
			// Register each enum value as a constant: EnumName.ValueName
			enumValue := &object.EnumValue{
				EnumName:  node.Name.Value,
				ValueName: v.Name.Value,
				Value:     value,
			}
			env.SetConst(node.Name.Value+"."+v.Name.Value, enumValue)
		}
		// Register the enum type itself
		env.SetConst(node.Name.Value, enumObj)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}

	case *ast.CharLiteral:
		return &object.Char{Value: node.Value}

	case *ast.StringLiteral:
		return evalStringLiteral(node.Value, env)

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.NullLiteral:
		return NULL

	case *ast.PrefixExpression:
		if node.Operator == "++" || node.Operator == "--" {
			return evalPrefixIncDec(node, env)
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		result := evalPrefixExpression(node.Operator, right)
		if errObj, ok := result.(*object.Error); ok && errObj.Line == 0 {
			errObj.Line = node.Token.Line
			errObj.Column = node.Token.Column
			errObj.EndColumn = node.Token.EndColumn
		}
		return result

	case *ast.PostfixExpression:
		return evalPostfixExpression(node, env)

	case *ast.InfixExpression:
		if node.Operator == "=" || node.Operator == "+=" || node.Operator == "-=" || node.Operator == "*=" || node.Operator == "/=" || node.Operator == "%=" {
			return evalAssignmentExpression(node, env)
		}

		if node.Operator == "." {
			return evalDotExpression(node, env)
		}

		// Short-circuit evaluation for && and ||
		if node.Operator == "&&" || node.Operator == "and" {
			left := Eval(node.Left, env)
			if isError(left) {
				return left
			}
			if !isTruthy(left) {
				return FALSE
			}
			right := Eval(node.Right, env)
			if isError(right) {
				return right
			}
			return nativeBoolToBooleanObject(isTruthy(right))
		}

		if node.Operator == "||" || node.Operator == "or" {
			left := Eval(node.Left, env)
			if isError(left) {
				return left
			}
			if isTruthy(left) {
				return TRUE
			}
			right := Eval(node.Right, env)
			if isError(right) {
				return right
			}
			return nativeBoolToBooleanObject(isTruthy(right))
		}

		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		result := evalInfixExpression(node.Operator, left, right)
		if errObj, ok := result.(*object.Error); ok && errObj.Line == 0 {
			errObj.Line = node.Token.Line
			errObj.Column = node.Token.Column
			errObj.EndColumn = node.Token.EndColumn
		}
		return result

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.TryStatement:
		return evalTryStatement(node, env)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{
			Parameters:      params,
			TypedParameters: node.TypedParameters,
			ReturnTypes:     node.ReturnTypes,
			Env:             env,
			Body:            body,
		}

	case *ast.ArrowFunction:
		params := node.Parameters
		body := node.Body
		return &object.ArrowFunction{Parameters: params, Env: env, Body: body}

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		result := applyFunction(function, args)
		if errObj, ok := result.(*object.Error); ok && errObj.Line == 0 {
			errObj.Line = node.Token.Line
			errObj.Column = node.Token.Column
			errObj.EndColumn = node.Token.EndColumn
		}
		return result

	case *ast.ArrayLiteral:
		elements := evalArrayElements(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		result := evalIndexExpression(left, index)
		if errObj, ok := result.(*object.Error); ok && errObj.Line == 0 {
			errObj.Line = node.Token.Line
			errObj.Column = node.Token.Column
			errObj.EndColumn = node.Token.EndColumn
		}
		return result

	case *ast.SliceExpression:
		return evalSliceExpression(node, env)

	case *ast.SpreadExpression:
		return newError("spread operator can only be used in array literals")

	case *ast.HashLiteral:
		return evalHashLiteral(node, env)

	case *ast.StructLiteral:
		s := &object.Struct{Name: node.Name.Value}
		for _, f := range node.Fields {
			s.Fields = append(s.Fields, f.Value)
		}
		env.Set(node.Name.Value, s)
		return NULL

	case *ast.StructInstantiation:
		return evalStructInstantiation(node, env)

	case *ast.MethodDefinition:
		fn := &object.Function{Parameters: node.Parameters, Env: env, Body: node.Body}
		key := node.StructName.Value + "." + node.MethodName.Value
		env.Set(key, fn)
		return NULL

	case *ast.WhileExpression:
		return evalWhileExpression(node, env)

	case *ast.ForExpression:
		return evalForExpression(node, env)

	case *ast.ForInIndexExpression:
		return evalForInIndexExpression(node, env)

	case *ast.CForExpression:
		return evalCForExpression(node, env)

	case *ast.BreakStatement:
		return &object.Break{}

	case *ast.ContinueStatement:
		return &object.Continue{}

	case *ast.SwitchExpression:
		return evalSwitchExpression(node, env)

	case *ast.TernaryExpression:
		return evalTernaryExpression(node, env)

	case *ast.RangeExpression:
		return evalRangeExpression(node, env)
	}

	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	blockEnv := object.NewEnclosedEnvironment(env)

	for _, statement := range block.Statements {
		result = Eval(statement, blockEnv)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ || rt == object.BREAK_OBJ || rt == object.CONTINUE_OBJ {
				return result
			}
		}
	}

	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newErrorWithLocation("identifier not found: "+node.Value, node.Token.Line, node.Token.Column, node.Token.EndColumn)
}

func evalStringLiteral(s string, env *object.Environment) object.Object {
	s = processEscapeSequences(s)

	result := ""
	i := 0
	for i < len(s) {
		if i+1 < len(s) && s[i] == '$' && s[i+1] == '{' {
			j := i + 2
			depth := 1
			for j < len(s) && depth > 0 {
				if s[j] == '{' {
					depth++
				} else if s[j] == '}' {
					depth--
				}
				j++
			}
			if depth == 0 {
				exprStr := s[i+2 : j-1]
				l := lexer.New(exprStr)
				p := parser.New(l)
				program := p.ParseProgram()
				if len(p.Errors()) > 0 {
					return newError("string interpolation parse error: %s", p.Errors()[0])
				}
				if len(program.Statements) == 0 {
					result += ""
				} else {
					val := Eval(program.Statements[0], env)
					if isError(val) {
						return val
					}
					if val != nil {
						result += val.Inspect()
					}
				}
				i = j
			} else {
				result += string(s[i])
				i++
			}
		} else {
			result += string(s[i])
			i++
		}
	}
	return &object.String{Value: result}
}

func processEscapeSequences(s string) string {
	result := ""
	i := 0
	for i < len(s) {
		if s[i] == '\\' && i+1 < len(s) {
			switch s[i+1] {
			case 'n':
				result += "\n"
				i += 2
			case 't':
				result += "\t"
				i += 2
			case 'r':
				result += "\r"
				i += 2
			case '\\':
				result += "\\"
				i += 2
			case '"':
				result += "\""
				i += 2
			case '$':
				result += "$"
				i += 2
			default:
				result += string(s[i])
				i++
			}
		} else {
			result += string(s[i])
			i++
		}
	}
	return result
}
