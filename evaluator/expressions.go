package evaluator

import (
	"victoria/ast"
	"victoria/object"
)

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() == object.INTEGER_OBJ {
		value := right.(*object.Integer).Value
		return &object.Integer{Value: -value}
	}
	if right.Type() == object.FLOAT_OBJ {
		value := right.(*object.Float).Value
		return &object.Float{Value: -value}
	}
	return newError("unknown operator: -%s", right.Type())
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(operator, left, right)
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(operator, &object.Float{Value: float64(left.(*object.Integer).Value)}, right)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalFloatInfixExpression(operator, left, &object.Float{Value: float64(right.(*object.Integer).Value)})
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return newError("division by zero")
		}
		return &object.Integer{Value: leftVal / rightVal}
	case "%":
		if rightVal == 0 {
			return newError("division by zero")
		}
		return &object.Integer{Value: leftVal % rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalFloatInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value

	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return newError("division by zero")
		}
		return &object.Float{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func evalTryStatement(node *ast.TryStatement, env *object.Environment) object.Object {
	result := Eval(node.Block, env)

	if isError(result) {
		if node.CatchBlock != nil {
			catchEnv := object.NewEnclosedEnvironment(env)
			if node.CatchVar != nil {
				msg := result.(*object.Error).Message
				catchEnv.Set(node.CatchVar.Value, &object.String{Value: msg})
			}
			return Eval(node.CatchBlock, catchEnv)
		}
		return NULL
	}

	return result
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func evalArrayElements(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		if spread, ok := e.(*ast.SpreadExpression); ok {
			evaluated := Eval(spread.Right, env)
			if isError(evaluated) {
				return []object.Object{evaluated}
			}
			if arr, ok := evaluated.(*object.Array); ok {
				result = append(result, arr.Elements...)
			} else {
				return []object.Object{newError("spread operator requires an array, got %s", evaluated.Type())}
			}
		} else {
			evaluated := Eval(e, env)
			if isError(evaluated) {
				return []object.Object{evaluated}
			}
			result = append(result, evaluated)
		}
	}

	return result
}

func evalSliceExpression(node *ast.SliceExpression, env *object.Environment) object.Object {
	left := Eval(node.Left, env)
	if isError(left) {
		return left
	}

	var startIdx, endIdx int64

	switch obj := left.(type) {
	case *object.Array:
		length := int64(len(obj.Elements))

		if node.Start != nil {
			startVal := Eval(node.Start, env)
			if isError(startVal) {
				return startVal
			}
			if intVal, ok := startVal.(*object.Integer); ok {
				startIdx = intVal.Value
				if startIdx < 0 {
					startIdx = length + startIdx
				}
			} else {
				return newError("slice index must be an integer, got %s", startVal.Type())
			}
		} else {
			startIdx = 0
		}

		if node.End != nil {
			endVal := Eval(node.End, env)
			if isError(endVal) {
				return endVal
			}
			if intVal, ok := endVal.(*object.Integer); ok {
				endIdx = intVal.Value
				if endIdx < 0 {
					endIdx = length + endIdx
				}
			} else {
				return newError("slice index must be an integer, got %s", endVal.Type())
			}
		} else {
			endIdx = length
		}

		if startIdx < 0 {
			startIdx = 0
		}
		if endIdx > length {
			endIdx = length
		}
		if startIdx > endIdx {
			return &object.Array{Elements: []object.Object{}}
		}

		newElements := make([]object.Object, endIdx-startIdx)
		copy(newElements, obj.Elements[startIdx:endIdx])
		return &object.Array{Elements: newElements}

	case *object.String:
		length := int64(len(obj.Value))

		if node.Start != nil {
			startVal := Eval(node.Start, env)
			if isError(startVal) {
				return startVal
			}
			if intVal, ok := startVal.(*object.Integer); ok {
				startIdx = intVal.Value
				if startIdx < 0 {
					startIdx = length + startIdx
				}
			} else {
				return newError("slice index must be an integer, got %s", startVal.Type())
			}
		} else {
			startIdx = 0
		}

		if node.End != nil {
			endVal := Eval(node.End, env)
			if isError(endVal) {
				return endVal
			}
			if intVal, ok := endVal.(*object.Integer); ok {
				endIdx = intVal.Value
				if endIdx < 0 {
					endIdx = length + endIdx
				}
			} else {
				return newError("slice index must be an integer, got %s", endVal.Type())
			}
		} else {
			endIdx = length
		}

		if startIdx < 0 {
			startIdx = 0
		}
		if endIdx > length {
			endIdx = length
		}
		if startIdx > endIdx {
			return &object.String{Value: ""}
		}

		return &object.String{Value: obj.Value[startIdx:endIdx]}

	default:
		return newError("slice operator not supported for: %s", left.Type())
	}
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return arrayObject.Elements[idx]
}

func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

func evalStructInstantiation(node *ast.StructInstantiation, env *object.Environment) object.Object {
	sObj, ok := env.Get(node.Name.Value)
	if !ok {
		return newError("struct not found: %s", node.Name.Value)
	}

	sDef, ok := sObj.(*object.Struct)
	if !ok {
		return newError("not a struct: %s", node.Name.Value)
	}

	instance := &object.StructInstance{Struct: sDef, Fields: make(map[string]object.Object)}

	for fieldName, expr := range node.Fields {
		val := Eval(expr, env)
		if isError(val) {
			return val
		}
		instance.Fields[fieldName] = val
	}

	return instance
}

func evalDotExpression(node *ast.InfixExpression, env *object.Environment) object.Object {
	left := Eval(node.Left, env)
	if isError(left) {
		return left
	}

	ident, ok := node.Right.(*ast.Identifier)
	if !ok {
		return newError("expected identifier after dot")
	}

	if left.Type() == object.HASH_OBJ {
		hash := left.(*object.Hash)
		key := &object.String{Value: ident.Value}
		hashed := key.HashKey()
		pair, ok := hash.Pairs[hashed]
		if ok {
			return pair.Value
		}
		return newError("property not found in hash: %s", ident.Value)
	}

	if left.Type() == object.INSTANCE_OBJ {
		instance := left.(*object.StructInstance)

		if val, ok := instance.Fields[ident.Value]; ok {
			return val
		}

		methodName := instance.Struct.Name + "." + ident.Value
		if method, ok := env.Get(methodName); ok {
			fn := method.(*object.Function)
			closureEnv := object.NewEnclosedEnvironment(fn.Env)
			closureEnv.Set("self", instance)

			return &object.Function{Parameters: fn.Parameters, Env: closureEnv, Body: fn.Body}
		}

		return newError("property or method not found: %s", ident.Value)
	}

	return newError("dot operator not supported for: %s", left.Type())
}

func evalPostfixExpression(node *ast.PostfixExpression, env *object.Environment) object.Object {
	ident, ok := node.Left.(*ast.Identifier)
	if !ok {
		return newError("postfix operator on non-identifier")
	}

	if env.IsConst(ident.Value) {
		return newErrorWithLocation("cannot reassign constant variable: "+ident.Value, node.Token.Line, node.Token.Column, node.Token.EndColumn)
	}

	currentVal, ok := env.Get(ident.Value)
	if !ok {
		return newError("variable not defined: %s", ident.Value)
	}

	var newVal object.Object
	one := &object.Integer{Value: 1}

	switch node.Operator {
	case "++":
		newVal = evalInfixExpression("+", currentVal, one)
	case "--":
		newVal = evalInfixExpression("-", currentVal, one)
	default:
		return newError("unknown operator: %s", node.Operator)
	}

	if isError(newVal) {
		return newVal
	}

	env.Update(ident.Value, newVal)
	return currentVal
}

func evalPrefixIncDec(node *ast.PrefixExpression, env *object.Environment) object.Object {
	ident, ok := node.Right.(*ast.Identifier)
	if !ok {
		return newError("prefix %s operator on non-identifier", node.Operator)
	}

	if env.IsConst(ident.Value) {
		return newErrorWithLocation("cannot reassign constant variable: "+ident.Value, node.Token.Line, node.Token.Column, node.Token.EndColumn)
	}

	currentVal, ok := env.Get(ident.Value)
	if !ok {
		return newError("variable not defined: %s", ident.Value)
	}

	var newVal object.Object
	one := &object.Integer{Value: 1}

	switch node.Operator {
	case "++":
		newVal = evalInfixExpression("+", currentVal, one)
	case "--":
		newVal = evalInfixExpression("-", currentVal, one)
	default:
		return newError("unknown operator: %s", node.Operator)
	}

	if isError(newVal) {
		return newVal
	}

	env.Update(ident.Value, newVal)
	return newVal
}

func evalAssignmentExpression(node *ast.InfixExpression, env *object.Environment) object.Object {
	ident, ok := node.Left.(*ast.Identifier)
	if !ok {
		return newError("assignment to non-identifier")
	}

	if env.IsConst(ident.Value) {
		return newErrorWithLocation("cannot reassign constant variable: "+ident.Value, node.Token.Line, node.Token.Column, node.Token.EndColumn)
	}

	if node.Operator == "=" {
		val := Eval(node.Right, env)
		if isError(val) {
			return val
		}

		_, ok = env.Update(ident.Value, val)
		if !ok {
			return newError("variable not defined: %s", ident.Value)
		}

		return val
	}

	currentVal, ok := env.Get(ident.Value)
	if !ok {
		return newError("variable not defined: %s", ident.Value)
	}

	rightVal := Eval(node.Right, env)
	if isError(rightVal) {
		return rightVal
	}

	var newVal object.Object
	switch node.Operator {
	case "+=":
		newVal = evalInfixExpression("+", currentVal, rightVal)
	case "-=":
		newVal = evalInfixExpression("-", currentVal, rightVal)
	case "*=":
		newVal = evalInfixExpression("*", currentVal, rightVal)
	case "/=":
		newVal = evalInfixExpression("/", currentVal, rightVal)
	case "%=":
		newVal = evalInfixExpression("%", currentVal, rightVal)
	}

	if isError(newVal) {
		return newVal
	}

	env.Update(ident.Value, newVal)
	return newVal
}

func evalTernaryExpression(node *ast.TernaryExpression, env *object.Environment) object.Object {
	condition := Eval(node.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(node.Consequence, env)
	}
	return Eval(node.Alternative, env)
}

func evalRangeExpression(node *ast.RangeExpression, env *object.Environment) object.Object {
	start := Eval(node.Start, env)
	if isError(start) {
		return start
	}

	end := Eval(node.End, env)
	if isError(end) {
		return end
	}

	startInt, ok := start.(*object.Integer)
	if !ok {
		return newError("range start must be an integer, got %s", start.Type())
	}

	endInt, ok := end.(*object.Integer)
	if !ok {
		return newError("range end must be an integer, got %s", end.Type())
	}

	return &object.Range{Start: startInt.Value, End: endInt.Value}
}
