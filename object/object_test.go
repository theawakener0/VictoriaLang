package object

import "testing"

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}
	diff1 := &String{Value: "My name is johnny"}
	diff2 := &String{Value: "My name is johnny"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("strings with different content have same hash keys")
	}
}

func TestIntegerHashKey(t *testing.T) {
	one1 := &Integer{Value: 1}
	one2 := &Integer{Value: 1}
	two1 := &Integer{Value: 2}
	two2 := &Integer{Value: 2}

	if one1.HashKey() != one2.HashKey() {
		t.Errorf("integers with same value have different hash keys")
	}

	if two1.HashKey() != two2.HashKey() {
		t.Errorf("integers with same value have different hash keys")
	}

	if one1.HashKey() == two1.HashKey() {
		t.Errorf("integers with different values have same hash keys")
	}
}

func TestBooleanHashKey(t *testing.T) {
	true1 := &Boolean{Value: true}
	true2 := &Boolean{Value: true}
	false1 := &Boolean{Value: false}
	false2 := &Boolean{Value: false}

	if true1.HashKey() != true2.HashKey() {
		t.Errorf("booleans with same value have different hash keys")
	}

	if false1.HashKey() != false2.HashKey() {
		t.Errorf("booleans with same value have different hash keys")
	}

	if true1.HashKey() == false1.HashKey() {
		t.Errorf("booleans with different values have same hash keys")
	}
}

func TestObjectTypes(t *testing.T) {
	tests := []struct {
		obj      Object
		expected ObjectType
	}{
		{&Integer{Value: 5}, INTEGER_OBJ},
		{&Float{Value: 3.14}, FLOAT_OBJ},
		{&String{Value: "hello"}, STRING_OBJ},
		{&Boolean{Value: true}, BOOLEAN_OBJ},
		{&Null{}, NULL_OBJ},
		{&Array{Elements: []Object{}}, ARRAY_OBJ},
		{&Hash{Pairs: map[HashKey]HashPair{}}, HASH_OBJ},
		{&Error{Message: "error"}, ERROR_OBJ},
		{&ReturnValue{Value: &Integer{Value: 5}}, RETURN_VALUE_OBJ},
		{&Break{}, BREAK_OBJ},
		{&Continue{}, CONTINUE_OBJ},
	}

	for _, tt := range tests {
		if tt.obj.Type() != tt.expected {
			t.Errorf("wrong type. got=%s, want=%s", tt.obj.Type(), tt.expected)
		}
	}
}

func TestIntegerInspect(t *testing.T) {
	obj := &Integer{Value: 42}
	if obj.Inspect() != "42" {
		t.Errorf("Integer.Inspect() wrong. got=%s, want=42", obj.Inspect())
	}
}

func TestFloatInspect(t *testing.T) {
	obj := &Float{Value: 3.14}
	if obj.Inspect() != "3.14" {
		t.Errorf("Float.Inspect() wrong. got=%s, want=3.14", obj.Inspect())
	}
}

func TestStringInspect(t *testing.T) {
	obj := &String{Value: "hello world"}
	if obj.Inspect() != "hello world" {
		t.Errorf("String.Inspect() wrong. got=%s, want='hello world'", obj.Inspect())
	}
}

func TestBooleanInspect(t *testing.T) {
	trueObj := &Boolean{Value: true}
	if trueObj.Inspect() != "true" {
		t.Errorf("Boolean.Inspect() wrong. got=%s, want=true", trueObj.Inspect())
	}

	falseObj := &Boolean{Value: false}
	if falseObj.Inspect() != "false" {
		t.Errorf("Boolean.Inspect() wrong. got=%s, want=false", falseObj.Inspect())
	}
}

func TestNullInspect(t *testing.T) {
	obj := &Null{}
	if obj.Inspect() != "null" {
		t.Errorf("Null.Inspect() wrong. got=%s, want=null", obj.Inspect())
	}
}

func TestErrorInspect(t *testing.T) {
	obj := &Error{Message: "something went wrong"}
	expected := "ERROR: something went wrong"
	if obj.Inspect() != expected {
		t.Errorf("Error.Inspect() wrong. got=%s, want=%s", obj.Inspect(), expected)
	}
}

func TestArrayInspect(t *testing.T) {
	obj := &Array{
		Elements: []Object{
			&Integer{Value: 1},
			&Integer{Value: 2},
			&Integer{Value: 3},
		},
	}
	expected := "[1, 2, 3]"
	if obj.Inspect() != expected {
		t.Errorf("Array.Inspect() wrong. got=%s, want=%s", obj.Inspect(), expected)
	}
}

func TestEnvironment(t *testing.T) {
	env := NewEnvironment()

	// Test Set and Get
	val := &Integer{Value: 42}
	env.Set("x", val)

	result, ok := env.Get("x")
	if !ok {
		t.Error("expected to find 'x' in environment")
	}
	if result != val {
		t.Errorf("wrong value. got=%v, want=%v", result, val)
	}

	// Test not found
	_, ok = env.Get("y")
	if ok {
		t.Error("expected 'y' to not be found")
	}
}

func TestEnclosedEnvironment(t *testing.T) {
	outer := NewEnvironment()
	outer.Set("x", &Integer{Value: 1})

	inner := NewEnclosedEnvironment(outer)
	inner.Set("y", &Integer{Value: 2})

	// Inner can access outer
	val, ok := inner.Get("x")
	if !ok {
		t.Error("expected inner to find 'x' from outer")
	}
	if val.(*Integer).Value != 1 {
		t.Errorf("wrong value for x. got=%d, want=1", val.(*Integer).Value)
	}

	// Inner can access its own
	val, ok = inner.Get("y")
	if !ok {
		t.Error("expected inner to find 'y'")
	}
	if val.(*Integer).Value != 2 {
		t.Errorf("wrong value for y. got=%d, want=2", val.(*Integer).Value)
	}

	// Outer cannot access inner
	_, ok = outer.Get("y")
	if ok {
		t.Error("expected outer to not find 'y'")
	}
}

func TestConstantVariable(t *testing.T) {
	env := NewEnvironment()

	// Set a constant
	val := &Integer{Value: 42}
	env.SetConst("PI", val)

	// Should be able to get it
	result, ok := env.Get("PI")
	if !ok {
		t.Error("expected to find 'PI' in environment")
	}
	if result != val {
		t.Errorf("wrong value. got=%v, want=%v", result, val)
	}

	// Should be identified as constant
	if !env.IsConst("PI") {
		t.Error("expected PI to be identified as constant")
	}
}

func TestReAssign(t *testing.T) {
	env := NewEnvironment()

	// Set a variable
	env.Set("x", &Integer{Value: 1})

	// Update it using Update method
	result, ok := env.Update("x", &Integer{Value: 2})
	if !ok {
		t.Error("expected Update to succeed")
	}

	// Check new value
	val, _ := env.Get("x")
	if val.(*Integer).Value != 2 {
		t.Errorf("wrong value after update. got=%d, want=2", val.(*Integer).Value)
	}
	if result.(*Integer).Value != 2 {
		t.Errorf("wrong return value from Update. got=%d, want=2", result.(*Integer).Value)
	}

	// Try to update non-existent variable
	_, ok = env.Update("y", &Integer{Value: 1})
	if ok {
		t.Error("expected Update to fail for non-existent variable")
	}
}

func TestReturnValue(t *testing.T) {
	inner := &Integer{Value: 42}
	obj := &ReturnValue{Value: inner}

	if obj.Type() != RETURN_VALUE_OBJ {
		t.Errorf("wrong type. got=%s, want=%s", obj.Type(), RETURN_VALUE_OBJ)
	}

	if obj.Inspect() != "42" {
		t.Errorf("wrong inspect. got=%s, want=42", obj.Inspect())
	}
}

func TestHashPairs(t *testing.T) {
	key := &String{Value: "name"}
	value := &String{Value: "Victoria"}

	pairs := map[HashKey]HashPair{
		key.HashKey(): {Key: key, Value: value},
	}

	hash := &Hash{Pairs: pairs}

	if hash.Type() != HASH_OBJ {
		t.Errorf("wrong type. got=%s, want=%s", hash.Type(), HASH_OBJ)
	}
}
