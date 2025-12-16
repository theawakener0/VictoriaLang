package evaluator

import (
	"victoria/ast"
	"victoria/object"
)

func evalWhileExpression(node *ast.WhileExpression, env *object.Environment) object.Object {
	var result object.Object = NULL

	for {
		condition := Eval(node.Condition, env)
		if isError(condition) {
			return condition
		}

		if !isTruthy(condition) {
			break
		}

		result = Eval(node.Body, env)
		if result != nil {
			if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
				return result
			}
			if result.Type() == object.BREAK_OBJ {
				break
			}
			if result.Type() == object.CONTINUE_OBJ {
				continue
			}
		}
	}

	return result
}

func evalForExpression(node *ast.ForExpression, env *object.Environment) object.Object {
	iterable := Eval(node.Iterable, env)
	if isError(iterable) {
		return iterable
	}

	var elements []object.Object

	switch iterable := iterable.(type) {
	case *object.Array:
		elements = iterable.Elements
	case *object.String:
		for _, char := range iterable.Value {
			elements = append(elements, &object.String{Value: string(char)})
		}
	case *object.Hash:
		for _, pair := range iterable.Pairs {
			elements = append(elements, pair.Key)
		}
	case *object.Range:
		for i := iterable.Start; i < iterable.End; i++ {
			elements = append(elements, &object.Integer{Value: i})
		}
	default:
		return newError("not iterable: %s", iterable.Type())
	}

	var result object.Object = NULL

	for _, elem := range elements {
		loopEnv := object.NewEnclosedEnvironment(env)
		loopEnv.Set(node.Item.Value, elem)

		result = evalBlockStatement(node.Body, loopEnv)

		if result != nil {
			if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
				return result
			}
			if result.Type() == object.BREAK_OBJ {
				break
			}
			if result.Type() == object.CONTINUE_OBJ {
				continue
			}
		}
	}

	return result
}

func evalForInIndexExpression(node *ast.ForInIndexExpression, env *object.Environment) object.Object {
	iterable := Eval(node.Iterable, env)
	if isError(iterable) {
		return iterable
	}

	var result object.Object = NULL

	switch iterable := iterable.(type) {
	case *object.Array:
		for i, elem := range iterable.Elements {
			loopEnv := object.NewEnclosedEnvironment(env)
			loopEnv.Set(node.Index.Value, &object.Integer{Value: int64(i)})
			loopEnv.Set(node.Value.Value, elem)

			result = evalBlockStatement(node.Body, loopEnv)

			if result != nil {
				if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
					return result
				}
				if result.Type() == object.BREAK_OBJ {
					break
				}
				if result.Type() == object.CONTINUE_OBJ {
					continue
				}
			}
		}
	case *object.Hash:
		for _, pair := range iterable.Pairs {
			loopEnv := object.NewEnclosedEnvironment(env)
			loopEnv.Set(node.Index.Value, pair.Key)
			loopEnv.Set(node.Value.Value, pair.Value)

			result = evalBlockStatement(node.Body, loopEnv)

			if result != nil {
				if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
					return result
				}
				if result.Type() == object.BREAK_OBJ {
					break
				}
				if result.Type() == object.CONTINUE_OBJ {
					continue
				}
			}
		}
	case *object.String:
		for i, char := range iterable.Value {
			loopEnv := object.NewEnclosedEnvironment(env)
			loopEnv.Set(node.Index.Value, &object.Integer{Value: int64(i)})
			loopEnv.Set(node.Value.Value, &object.String{Value: string(char)})

			result = evalBlockStatement(node.Body, loopEnv)

			if result != nil {
				if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
					return result
				}
				if result.Type() == object.BREAK_OBJ {
					break
				}
				if result.Type() == object.CONTINUE_OBJ {
					continue
				}
			}
		}
	default:
		return newError("not iterable: %s", iterable.Type())
	}

	return result
}

func evalCForExpression(node *ast.CForExpression, env *object.Environment) object.Object {
	loopEnv := object.NewEnclosedEnvironment(env)

	if node.Init != nil {
		initResult := Eval(node.Init, loopEnv)
		if isError(initResult) {
			return initResult
		}
	}

	var result object.Object = NULL

	for {
		if node.Condition != nil {
			cond := Eval(node.Condition, loopEnv)
			if isError(cond) {
				return cond
			}
			if !isTruthy(cond) {
				break
			}
		}

		bodyEnv := object.NewEnclosedEnvironment(loopEnv)
		result = evalBlockStatement(node.Body, bodyEnv)

		if result != nil {
			if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
				return result
			}
			if result.Type() == object.BREAK_OBJ {
				break
			}
			if result.Type() == object.CONTINUE_OBJ {
				// Continue: fall through to update expression
				_ = result // Explicit no-op to satisfy linter
			}
		}

		if node.Update != nil {
			updateResult := Eval(node.Update, loopEnv)
			if isError(updateResult) {
				return updateResult
			}
		}
	}

	return result
}

func evalSwitchExpression(node *ast.SwitchExpression, env *object.Environment) object.Object {
	value := Eval(node.Value, env)
	if isError(value) {
		return value
	}

	for _, caseExpr := range node.Cases {
		caseValue := Eval(caseExpr.Value, env)
		if isError(caseValue) {
			return caseValue
		}

		if compareObjects(value, caseValue) {
			return Eval(caseExpr.Body, env)
		}
	}

	if node.Default != nil {
		return Eval(node.Default, env)
	}

	return NULL
}

func compareObjects(a, b object.Object) bool {
	if a.Type() != b.Type() {
		return false
	}
	switch a := a.(type) {
	case *object.Integer:
		return a.Value == b.(*object.Integer).Value
	case *object.String:
		return a.Value == b.(*object.String).Value
	case *object.Boolean:
		return a.Value == b.(*object.Boolean).Value
	case *object.Float:
		return a.Value == b.(*object.Float).Value
	}
	return a == b
}
