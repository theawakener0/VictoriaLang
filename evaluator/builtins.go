package evaluator

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"victoria/object"
)

// builtins is a map of built-in functions available in the Victoria language
var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"print": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},
	"range": {
		Fn: func(args ...object.Object) object.Object {
			var start, end, step int64

			switch len(args) {
			case 1:
				endVal, ok := args[0].(*object.Integer)
				if !ok {
					return newError("argument to `range` must be INTEGER, got %s", args[0].Type())
				}
				start, end, step = 0, endVal.Value, 1
			case 2:
				startVal, ok1 := args[0].(*object.Integer)
				endVal, ok2 := args[1].(*object.Integer)
				if !ok1 || !ok2 {
					return newError("arguments to `range` must be integers")
				}
				start, end, step = startVal.Value, endVal.Value, 1
			case 3:
				startVal, ok1 := args[0].(*object.Integer)
				endVal, ok2 := args[1].(*object.Integer)
				stepVal, ok3 := args[2].(*object.Integer)
				if !ok1 || !ok2 || !ok3 {
					return newError("arguments to `range` must be integers")
				}
				if stepVal.Value == 0 {
					return newError("range step cannot be zero")
				}
				start, end, step = startVal.Value, endVal.Value, stepVal.Value
			default:
				return newError("wrong number of arguments. got=%d, want=1, 2, or 3", len(args))
			}

			elements := []object.Object{}
			if step > 0 {
				for i := start; i < end; i += step {
					elements = append(elements, &object.Integer{Value: i})
				}
			} else {
				for i := start; i > end; i += step {
					elements = append(elements, &object.Integer{Value: i})
				}
			}

			return &object.Array{Elements: elements}
		},
	},
	"format": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 1 {
				return newError("wrong number of arguments. got=%d, want=at least 1", len(args))
			}

			formatStr, ok := args[0].(*object.String)
			if !ok {
				return newError("argument 1 to `format` must be STRING, got %s", args[0].Type())
			}

			var fmtArgs []interface{}
			for _, arg := range args[1:] {
				fmtArgs = append(fmtArgs, unwrapObject(arg))
			}

			return &object.String{Value: fmt.Sprintf(formatStr.Value, fmtArgs...)}
		},
	},
	"input": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) > 1 {
				return newError("wrong number of arguments. got=%d, want=0 or 1", len(args))
			}

			if len(args) == 1 {
				fmt.Print(args[0].Inspect())
			}

			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			text = strings.TrimSpace(text)

			return &object.String{Value: text}
		},
	},
	"int": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Integer:
				return arg
			case *object.String:
				val, err := strconv.ParseInt(arg.Value, 0, 64)
				if err != nil {
					return newError("could not parse %q as integer", arg.Value)
				}
				return &object.Integer{Value: val}
			case *object.Boolean:
				if arg.Value {
					return &object.Integer{Value: 1}
				}
				return &object.Integer{Value: 0}
			default:
				return newError("argument to `int` not supported, got %s", args[0].Type())
			}
		},
	},
	"char": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Integer:
				return &object.Char{Value: rune(arg.Value)}
			case *object.String:
				if len(arg.Value) == 0 {
					return newError("cannot convert empty string to char")
				}
				runes := []rune(arg.Value)
				return &object.Char{Value: runes[0]}
			case *object.Char:
				return arg
			case *object.Byte:
				return &object.Char{Value: rune(arg.Value)}
			case *object.Rune:
				return &object.Char{Value: arg.Value}
			default:
				return newError("argument to `char` not supported, got %s", args[0].Type())
			}
		},
	},
	"byte": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Integer:
				return &object.Byte{Value: byte(arg.Value)}
			case *object.String:
				if len(arg.Value) == 0 {
					return newError("cannot convert empty string to byte")
				}
				return &object.Byte{Value: arg.Value[0]}
			case *object.Char:
				return &object.Byte{Value: byte(arg.Value)}
			case *object.Byte:
				return arg
			case *object.Rune:
				return &object.Byte{Value: byte(arg.Value)}
			default:
				return newError("argument to `byte` not supported, got %s", args[0].Type())
			}
		},
	},
	"rune": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Integer:
				return &object.Rune{Value: rune(arg.Value)}
			case *object.String:
				if len(arg.Value) == 0 {
					return newError("cannot convert empty string to rune")
				}
				runes := []rune(arg.Value)
				return &object.Rune{Value: runes[0]}
			case *object.Char:
				return &object.Rune{Value: arg.Value}
			case *object.Byte:
				return &object.Rune{Value: rune(arg.Value)}
			case *object.Rune:
				return arg
			default:
				return newError("argument to `rune` not supported, got %s", args[0].Type())
			}
		},
	},
	"ord": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Char:
				return &object.Integer{Value: int64(arg.Value)}
			case *object.String:
				if len(arg.Value) == 0 {
					return newError("cannot get ordinal of empty string")
				}
				runes := []rune(arg.Value)
				return &object.Integer{Value: int64(runes[0])}
			case *object.Byte:
				return &object.Integer{Value: int64(arg.Value)}
			case *object.Rune:
				return &object.Integer{Value: int64(arg.Value)}
			default:
				return newError("argument to `ord` not supported, got %s", args[0].Type())
			}
		},
	},
	"chr": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Integer:
				return &object.String{Value: string(rune(arg.Value))}
			case *object.Char:
				return &object.String{Value: string(arg.Value)}
			case *object.Byte:
				return &object.String{Value: string(arg.Value)}
			case *object.Rune:
				return &object.String{Value: string(arg.Value)}
			default:
				return newError("argument to `chr` not supported, got %s", args[0].Type())
			}
		},
	},
	"isDigit": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			var r rune
			switch arg := args[0].(type) {
			case *object.Char:
				r = arg.Value
			case *object.String:
				if len(arg.Value) == 0 {
					return FALSE
				}
				runes := []rune(arg.Value)
				r = runes[0]
			case *object.Rune:
				r = arg.Value
			default:
				return newError("argument to `isDigit` not supported, got %s", args[0].Type())
			}
			if r >= '0' && r <= '9' {
				return TRUE
			}
			return FALSE
		},
	},
	"isLetter": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			var r rune
			switch arg := args[0].(type) {
			case *object.Char:
				r = arg.Value
			case *object.String:
				if len(arg.Value) == 0 {
					return FALSE
				}
				runes := []rune(arg.Value)
				r = runes[0]
			case *object.Rune:
				r = arg.Value
			default:
				return newError("argument to `isLetter` not supported, got %s", args[0].Type())
			}
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				return TRUE
			}
			return FALSE
		},
	},
	"isAlpha": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			var r rune
			switch arg := args[0].(type) {
			case *object.Char:
				r = arg.Value
			case *object.String:
				if len(arg.Value) == 0 {
					return FALSE
				}
				runes := []rune(arg.Value)
				r = runes[0]
			case *object.Rune:
				r = arg.Value
			default:
				return newError("argument to `isAlpha` not supported, got %s", args[0].Type())
			}
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
				return TRUE
			}
			return FALSE
		},
	},
	"isSpace": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			var r rune
			switch arg := args[0].(type) {
			case *object.Char:
				r = arg.Value
			case *object.String:
				if len(arg.Value) == 0 {
					return FALSE
				}
				runes := []rune(arg.Value)
				r = runes[0]
			case *object.Rune:
				r = arg.Value
			default:
				return newError("argument to `isSpace` not supported, got %s", args[0].Type())
			}
			if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
				return TRUE
			}
			return FALSE
		},
	},
	"toUpper": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Char:
				r := arg.Value
				if r >= 'a' && r <= 'z' {
					return &object.Char{Value: r - 32}
				}
				return arg
			case *object.Rune:
				r := arg.Value
				if r >= 'a' && r <= 'z' {
					return &object.Rune{Value: r - 32}
				}
				return arg
			case *object.String:
				return &object.String{Value: strings.ToUpper(arg.Value)}
			default:
				return newError("argument to `toUpper` not supported, got %s", args[0].Type())
			}
		},
	},
	"toLower": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Char:
				r := arg.Value
				if r >= 'A' && r <= 'Z' {
					return &object.Char{Value: r + 32}
				}
				return arg
			case *object.Rune:
				r := arg.Value
				if r >= 'A' && r <= 'Z' {
					return &object.Rune{Value: r + 32}
				}
				return arg
			case *object.String:
				return &object.String{Value: strings.ToLower(arg.Value)}
			default:
				return newError("argument to `toLower` not supported, got %s", args[0].Type())
			}
		},
	},
	"string": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			return &object.String{Value: args[0].Inspect()}
		},
	},
	"type": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			return &object.String{Value: string(args[0].Type())}
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}
			return NULL
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `last` must be ARRAY, got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}
			return NULL
		},
	},
	"rest": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `rest` must be ARRAY, got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				newElements := make([]object.Object, length-1)
				copy(newElements, arr.Elements[1:length])
				return &object.Array{Elements: newElements}
			}
			return NULL
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `push` must be ARRAY, got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			newElements := make([]object.Object, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]
			return &object.Array{Elements: newElements}
		},
	},
	"pop": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `pop` must be ARRAY, got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				newElements := make([]object.Object, length-1)
				copy(newElements, arr.Elements[:length-1])
				return &object.Array{Elements: newElements}
			}
			return NULL
		},
	},
	"split": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != object.STRING_OBJ {
				return newError("argument 1 to `split` must be STRING, got %s", args[0].Type())
			}
			if args[1].Type() != object.STRING_OBJ {
				return newError("argument 2 to `split` must be STRING, got %s", args[1].Type())
			}
			str := args[0].(*object.String).Value
			sep := args[1].(*object.String).Value
			parts := strings.Split(str, sep)
			elements := make([]object.Object, len(parts))
			for i, p := range parts {
				elements[i] = &object.String{Value: p}
			}
			return &object.Array{Elements: elements}
		},
	},
	"join": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument 1 to `join` must be ARRAY, got %s", args[0].Type())
			}
			if args[1].Type() != object.STRING_OBJ {
				return newError("argument 2 to `join` must be STRING, got %s", args[1].Type())
			}
			arr := args[0].(*object.Array)
			sep := args[1].(*object.String).Value
			parts := make([]string, len(arr.Elements))
			for i, e := range arr.Elements {
				if e.Type() != object.STRING_OBJ {
					return newError("array elements must be strings")
				}
				parts[i] = e.(*object.String).Value
			}
			return &object.String{Value: strings.Join(parts, sep)}
		},
	},
	"contains": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			switch container := args[0].(type) {
			case *object.Array:
				for _, e := range container.Elements {
					if e.Type() == args[1].Type() {
						if e.Inspect() == args[1].Inspect() {
							return TRUE
						}
					}
				}
				return FALSE
			case *object.String:
				if args[1].Type() != object.STRING_OBJ {
					return newError("argument 2 to `contains` on string must be STRING")
				}
				if strings.Contains(container.Value, args[1].(*object.String).Value) {
					return TRUE
				}
				return FALSE
			default:
				return newError("argument 1 to `contains` must be ARRAY or STRING, got %s", args[0].Type())
			}
		},
	},
	"index": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			switch container := args[0].(type) {
			case *object.Array:
				for i, e := range container.Elements {
					if e.Type() == args[1].Type() && e.Inspect() == args[1].Inspect() {
						return &object.Integer{Value: int64(i)}
					}
				}
				return &object.Integer{Value: -1}
			case *object.String:
				if args[1].Type() != object.STRING_OBJ {
					return newError("argument 2 to `index` on string must be STRING")
				}
				idx := strings.Index(container.Value, args[1].(*object.String).Value)
				return &object.Integer{Value: int64(idx)}
			default:
				return newError("argument 1 to `index` must be ARRAY or STRING, got %s", args[0].Type())
			}
		},
	},
	"upper": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.STRING_OBJ {
				return newError("argument to `upper` must be STRING, got %s", args[0].Type())
			}
			return &object.String{Value: strings.ToUpper(args[0].(*object.String).Value)}
		},
	},
	"lower": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.STRING_OBJ {
				return newError("argument to `lower` must be STRING, got %s", args[0].Type())
			}
			return &object.String{Value: strings.ToLower(args[0].(*object.String).Value)}
		},
	},
	"keys": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.HASH_OBJ {
				return newError("argument to `keys` must be HASH, got %s", args[0].Type())
			}
			hash := args[0].(*object.Hash)
			elements := []object.Object{}
			for _, pair := range hash.Pairs {
				elements = append(elements, pair.Key)
			}
			return &object.Array{Elements: elements}
		},
	},
	"values": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.HASH_OBJ {
				return newError("argument to `values` must be HASH, got %s", args[0].Type())
			}
			hash := args[0].(*object.Hash)
			elements := []object.Object{}
			for _, pair := range hash.Pairs {
				elements = append(elements, pair.Value)
			}
			return &object.Array{Elements: elements}
		},
	},
	"map":    nil, // initialized in init()
	"filter": nil, // initialized in init()
	"reduce": nil, // initialized in init()
}

func init() {
	builtins["map"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument 1 to `map` must be ARRAY, got %s", args[0].Type())
			}
			if !isCallable(args[1]) {
				return newError("argument 2 to `map` must be FUNCTION, got %s", args[1].Type())
			}
			arr := args[0].(*object.Array)
			fn := args[1]
			paramCount := getParamCount(fn)
			elements := make([]object.Object, len(arr.Elements))
			for i, e := range arr.Elements {
				fnArgs := []object.Object{e}
				if paramCount > 1 {
					fnArgs = append(fnArgs, &object.Integer{Value: int64(i)})
				}
				result := applyFunction(fn, fnArgs)
				if isError(result) {
					return result
				}
				elements[i] = result
			}
			return &object.Array{Elements: elements}
		},
	}
	builtins["filter"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument 1 to `filter` must be ARRAY, got %s", args[0].Type())
			}
			if !isCallable(args[1]) {
				return newError("argument 2 to `filter` must be FUNCTION, got %s", args[1].Type())
			}
			arr := args[0].(*object.Array)
			fn := args[1]
			paramCount := getParamCount(fn)
			elements := []object.Object{}
			for i, e := range arr.Elements {
				fnArgs := []object.Object{e}
				if paramCount > 1 {
					fnArgs = append(fnArgs, &object.Integer{Value: int64(i)})
				}
				result := applyFunction(fn, fnArgs)
				if isError(result) {
					return result
				}
				if isTruthy(result) {
					elements = append(elements, e)
				}
			}
			return &object.Array{Elements: elements}
		},
	}
	builtins["reduce"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 2 || len(args) > 3 {
				return newError("wrong number of arguments. got=%d, want=2 or 3", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument 1 to `reduce` must be ARRAY, got %s", args[0].Type())
			}
			if !isCallable(args[1]) {
				return newError("argument 2 to `reduce` must be FUNCTION, got %s", args[1].Type())
			}
			arr := args[0].(*object.Array)
			fn := args[1]
			paramCount := getParamCount(fn)

			var accumulator object.Object
			startIdx := 0

			if len(args) == 3 {
				accumulator = args[2]
			} else if len(arr.Elements) > 0 {
				accumulator = arr.Elements[0]
				startIdx = 1
			} else {
				return newError("reduce of empty array with no initial value")
			}

			for i := startIdx; i < len(arr.Elements); i++ {
				fnArgs := []object.Object{accumulator, arr.Elements[i]}
				if paramCount > 2 {
					fnArgs = append(fnArgs, &object.Integer{Value: int64(i)})
				}
				result := applyFunction(fn, fnArgs)
				if isError(result) {
					return result
				}
				accumulator = result
			}
			return accumulator
		},
	}
}
