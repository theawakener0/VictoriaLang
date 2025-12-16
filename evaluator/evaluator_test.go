package evaluator

import (
	"testing"
	"victoria/lexer"
	"victoria/object"
	"victoria/parser"
)

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	RegisterBuiltinModules()
	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}
	return true
}

func testFloatObject(t *testing.T, obj object.Object, expected float64) bool {
	result, ok := obj.(*object.Float)
	if !ok {
		t.Errorf("object is not Float. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%f, want=%f", result.Value, expected)
		return false
	}
	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}
	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
		{"10 % 3", 1},
		{"17 % 5", 2},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalFloatExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"3.14", 3.14},
		{"0.5", 0.5},
		{"-3.14", -3.14},
		{"1.5 + 2.5", 4.0},
		{"5.0 - 2.0", 3.0},
		{"2.0 * 3.0", 6.0},
		{"10.0 / 4.0", 2.5},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testFloatObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{"1 <= 1", true},
		{"1 >= 1", true},
		{"1 <= 2", true},
		{"2 >= 1", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{`if (10 > 1) {
			if (10 > 1) {
				return 10;
			}
			return 1;
		}`, 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{"5 + true;", "type mismatch: INTEGER + BOOLEAN"},
		{"5 + true; 5;", "type mismatch: INTEGER + BOOLEAN"},
		{"-true", "unknown operator: -BOOLEAN"},
		{"true + false;", "unknown operator: BOOLEAN + BOOLEAN"},
		{"5; true + false; 5", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) { true + false; }", "unknown operator: BOOLEAN + BOOLEAN"},
		{"foobar", "identifier not found: foobar"},
		{`"Hello" - "World"`, "unknown operator: STRING - STRING"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "let identity = define(x) { x + 2 }; identity;"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = define(x) { x }; identity(5)", 5},
		{"let identity = define(x) { return x }; identity(5)", 5},
		{"let double = define(x) { x * 2 }; double(5)", 10},
		{"let add = define(x, y) { x + y }; add(5, 5)", 10},
		{"let add = define(x, y) { x + y }; add(5 + 5, add(5, 5))", 20},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
define newAdder(x) {
    return define(y) { x + y }
}

let addTwo = newAdder(2)
addTwo(2)`

	testIntegerObject(t, testEval(input), 4)
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
		{`len([1, 2, 3])`, 3},
		{`len([])`, 0},
		{`first([1, 2, 3])`, 1},
		{`first([])`, nil},
		{`last([1, 2, 3])`, 3},
		{`last([])`, nil},
		{`rest([1, 2, 3])`, []int{2, 3}},
		{`push([], 1)`, []int{1}},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case nil:
			testNullObject(t, evaluated)
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q",
					expected, errObj.Message)
			}
		case []int:
			array, ok := evaluated.(*object.Array)
			if !ok {
				t.Errorf("obj not Array. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if len(array.Elements) != len(expected) {
				t.Errorf("wrong num of elements. want=%d, got=%d",
					len(expected), len(array.Elements))
				continue
			}
			for i, expectedElem := range expected {
				testIntegerObject(t, array.Elements[i], int64(expectedElem))
			}
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d", len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"[1, 2, 3][0]", 1},
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][2]", 3},
		{"let i = 0; [1][i];", 1},
		{"[1, 2, 3][1 + 1];", 3},
		{"let myArray = [1, 2, 3]; myArray[2];", 3},
		{"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];", 6},
		{"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]", 2},
		{"[1, 2, 3][3]", nil},
		{"[1, 2, 3][-1]", nil},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestHashLiterals(t *testing.T) {
	input := `let two = "two";
	{
		"one": 10 - 9,
		two: 1 + 1,
		"thr" + "ee": 6 / 2,
		4: 4,
		true: 5,
		false: 6
	}`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, evaluated)
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}

		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`{"foo": 5}["foo"]`, 5},
		{`{"foo": 5}["bar"]`, nil},
		{`let key = "foo"; {"foo": 5}[key]`, 5},
		{`{}["foo"]`, nil},
		{`{5: 5}[5]`, 5},
		{`{true: 5}[true]`, 5},
		{`{false: 5}[false]`, 5},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestDivisionByZero(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{"10 / 0", "division by zero"},
		{"10 % 0", "division by zero"},
		{"10.0 / 0.0", "division by zero"},
		{"let x = 0; 10 / x", "division by zero"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned for %q. got=%T(%+v)",
				tt.input, evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

func TestWhileLoops(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let x = 0; while (x < 5) { x = x + 1; }; x;", 5},
		{"let x = 10; while (x > 0) { x = x - 2; }; x;", 0},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestForLoops(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let sum = 0; for x in [1, 2, 3] { sum = sum + x; }; sum;", 6},
		{"let sum = 0; for i in 1..5 { sum = sum + i; }; sum;", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestArrowFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let double = x => x * 2; double(5);", 10},
		{"let add = (a, b) => a + b; add(3, 4);", 7},
		{"(x => x + 1)(5)", 6},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestTernaryExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"true ? 1 : 2", 1},
		{"false ? 1 : 2", 2},
		{"1 < 2 ? 10 : 20", 10},
		{"1 > 2 ? 10 : 20", 20},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestRangeExpression(t *testing.T) {
	// Test using range in a for loop, which is how it's typically used
	input := `
let sum = 0
for i in 1..5 {
    sum = sum + i
}
sum`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 10) // 1+2+3+4 = 10
}

func TestArraySlicing(t *testing.T) {
	tests := []struct {
		input    string
		expected []int64
	}{
		{"[1, 2, 3, 4, 5][1:3]", []int64{2, 3}},
		{"[1, 2, 3, 4, 5][:2]", []int64{1, 2}},
		{"[1, 2, 3, 4, 5][3:]", []int64{4, 5}},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		arr, ok := evaluated.(*object.Array)
		if !ok {
			t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
		}

		if len(arr.Elements) != len(tt.expected) {
			t.Fatalf("array has wrong length. got=%d, want=%d",
				len(arr.Elements), len(tt.expected))
		}

		for i, exp := range tt.expected {
			testIntegerObject(t, arr.Elements[i], exp)
		}
	}
}

func TestSpreadOperator(t *testing.T) {
	input := "let a = [1, 2]; let b = [0, ...a, 3]; b;"

	evaluated := testEval(input)
	arr, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}

	expected := []int64{0, 1, 2, 3}
	if len(arr.Elements) != len(expected) {
		t.Fatalf("array has wrong length. got=%d, want=%d",
			len(arr.Elements), len(expected))
	}

	for i, exp := range expected {
		testIntegerObject(t, arr.Elements[i], exp)
	}
}

func TestStringInterpolation(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`let x = 5; "value is ${x}"`, "value is 5"},
		{`let name = "Victoria"; "Hello, ${name}!"`, "Hello, Victoria!"},
		{`"1 + 2 = ${1 + 2}"`, "1 + 2 = 3"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		str, ok := evaluated.(*object.String)
		if !ok {
			t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
		}

		if str.Value != tt.expected {
			t.Errorf("String has wrong value. got=%q, want=%q", str.Value, tt.expected)
		}
	}
}

func TestConstStatement(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{"const X = 5; X = 10", "cannot reassign constant variable: X"},
		{"const Y = 10; Y += 5", "cannot reassign constant variable: Y"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned for %q. got=%T(%+v)",
				tt.input, evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

func TestLogicalOperators(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true and true", true},
		{"true and false", false},
		{"false and true", false},
		{"false and false", false},
		{"true or true", true},
		{"true or false", true},
		{"false or true", true},
		{"false or false", false},
		{"true && true", true},
		{"true && false", false},
		{"true || false", true},
		{"false || false", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestBreakStatement(t *testing.T) {
	input := `
let sum = 0
for i in 1..10 {
    if (i > 5) { break }
    sum = sum + i
}
sum`

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 15) // 1 + 2 + 3 + 4 + 5 = 15
}

func TestContinueStatement(t *testing.T) {
	input := `
let sum = 0
for i in 1..6 {
    if (i == 3) { continue }
    sum = sum + i
}
sum`

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 12) // 1 + 2 + 4 + 5 = 12
}

func TestTryCatch(t *testing.T) {
	input := `
let result = try {
    let x = undefined_var;
    10;
} catch(e) {
    20;
}
result;`

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 20)
}

func TestSwitchStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{`let x = 1; switch (x) { case 1: { 10 } case 2: { 20 } }`, 10},
		{`let x = 2; switch (x) { case 1: { 10 } case 2: { 20 } }`, 20},
		{`let x = 3; switch (x) { case 1: { 10 } default: { 30 } }`, 30},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestIncrementDecrement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let x = 5; x++; x;", 6},
		{"let x = 5; x--; x;", 4},
		{"let x = 5; ++x;", 6},
		{"let x = 5; --x;", 4},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestCompoundAssignment(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let x = 5; x += 3; x;", 8},
		{"let x = 10; x -= 3; x;", 7},
		{"let x = 5; x *= 3; x;", 15},
		{"let x = 15; x /= 3; x;", 5},
		{"let x = 17; x %= 5; x;", 2},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestHigherOrderFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// Map
		{`map([1, 2, 3], x => x * 2)`, []int64{2, 4, 6}},
		// Filter
		{`filter([1, 2, 3, 4], x => x > 2)`, []int64{3, 4}},
		// Reduce
		{`reduce([1, 2, 3, 4], (acc, x) => acc + x, 0)`, int64(10)},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case []int64:
			arr, ok := evaluated.(*object.Array)
			if !ok {
				t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
			}

			if len(arr.Elements) != len(expected) {
				t.Fatalf("array has wrong length. got=%d, want=%d",
					len(arr.Elements), len(expected))
			}

			for i, exp := range expected {
				testIntegerObject(t, arr.Elements[i], exp)
			}
		case int64:
			testIntegerObject(t, evaluated, expected)
		}
	}
}
