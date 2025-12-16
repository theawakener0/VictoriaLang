package evaluator

import (
	"victoria/object"
)

func applyFunction(fn object.Object, args []object.Object) object.Object {
	if fn == nil {
		return newError("not a function: nil")
	}
	switch fn := fn.(type) {
	case *object.Function:
		// Type check arguments if typed parameters are present
		if len(fn.TypedParameters) > 0 {
			if len(args) != len(fn.TypedParameters) {
				return newError("wrong number of arguments: expected %d, got %d",
					len(fn.TypedParameters), len(args))
			}
			for i, typedParam := range fn.TypedParameters {
				if typedParam.Type != nil {
					if !object.CheckType(args[i], typedParam.Type) {
						return newError("type mismatch for parameter '%s': expected %s, got %s",
							typedParam.Name.Value, typedParam.Type.String(), object.TypeName(args[i]))
					}
				}
			}
		}

		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		result := unwrapReturnValue(evaluated)

		// Type check return value if return types are specified
		if len(fn.ReturnTypes) > 0 && !isError(result) {
			// For single return type
			if len(fn.ReturnTypes) == 1 {
				if !object.CheckType(result, fn.ReturnTypes[0]) {
					return newError("return type mismatch: expected %s, got %s",
						fn.ReturnTypes[0].String(), object.TypeName(result))
				}
			}
			// Multiple return types could be supported with tuple unpacking
		}

		return result

	case *object.ArrowFunction:
		extendedEnv := extendArrowFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return fn.Fn(args...)

	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for i, param := range fn.Parameters {
		env.Set(param.Value, args[i])
	}
	return env
}

func extendArrowFunctionEnv(fn *object.ArrowFunction, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for i, param := range fn.Parameters {
		if i < len(args) {
			env.Set(param.Value, args[i])
		}
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

// isCallable checks if an object can be called as a function
func isCallable(obj object.Object) bool {
	switch obj.Type() {
	case object.FUNCTION_OBJ, object.ARROW_FUNCTION_OBJ, object.BUILTIN_OBJ:
		return true
	default:
		return false
	}
}

// getParamCount returns the number of parameters for a callable
func getParamCount(obj object.Object) int {
	switch fn := obj.(type) {
	case *object.Function:
		return len(fn.Parameters)
	case *object.ArrowFunction:
		return len(fn.Parameters)
	default:
		return 0
	}
}

// unwrapObject converts a Victoria object to a Go interface{}
func unwrapObject(obj object.Object) interface{} {
	switch obj := obj.(type) {
	case *object.Integer:
		return obj.Value
	case *object.Boolean:
		return obj.Value
	case *object.String:
		return obj.Value
	case *object.Null:
		return nil
	default:
		return obj.Inspect()
	}
}
