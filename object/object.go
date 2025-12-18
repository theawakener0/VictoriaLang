package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"
	"victoria/ast"
)

type ObjectType string

const (
	INTEGER_OBJ        = "INTEGER"
	FLOAT_OBJ          = "FLOAT"
	BOOLEAN_OBJ        = "BOOLEAN"
	NULL_OBJ           = "NULL"
	RETURN_VALUE_OBJ   = "RETURN_VALUE"
	ERROR_OBJ          = "ERROR"
	FUNCTION_OBJ       = "FUNCTION"
	ARROW_FUNCTION_OBJ = "ARROW_FUNCTION"
	STRING_OBJ         = "STRING"
	CHAR_OBJ           = "CHAR"
	BYTE_OBJ           = "BYTE"
	RUNE_OBJ           = "RUNE"
	BUILTIN_OBJ        = "BUILTIN"
	ARRAY_OBJ          = "ARRAY"
	HASH_OBJ           = "HASH"
	STRUCT_OBJ         = "STRUCT"          // The struct definition
	INSTANCE_OBJ       = "STRUCT_INSTANCE" // The instance
	ENUM_OBJ           = "ENUM"            // Enum type definition
	ENUM_VALUE_OBJ     = "ENUM_VALUE"      // Enum value
	BREAK_OBJ          = "BREAK"
	CONTINUE_OBJ       = "CONTINUE"
	RANGE_OBJ          = "RANGE"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Hashable interface {
	HashKey() HashKey
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

type Float struct {
	Value float64
}

func (f *Float) Type() ObjectType { return FLOAT_OBJ }
func (f *Float) Inspect() string  { return fmt.Sprintf("%g", f.Value) }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type Error struct {
	Message   string
	Line      int
	Column    int
	EndColumn int
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

type Function struct {
	Parameters      []*ast.Identifier
	TypedParameters []*ast.TypedParameter // Parameters with type annotations
	ReturnTypes     []*ast.TypeAnnotation // Return type(s)
	Body            *ast.BlockStatement
	Env             *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	if len(f.TypedParameters) > 0 {
		for _, p := range f.TypedParameters {
			params = append(params, p.String())
		}
	} else {
		for _, p := range f.Parameters {
			params = append(params, p.String())
		}
	}
	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	if len(f.ReturnTypes) > 0 {
		out.WriteString(" -> ")
		types := []string{}
		for _, rt := range f.ReturnTypes {
			types = append(types, rt.String())
		}
		out.WriteString(strings.Join(types, ", "))
	}
	out.WriteString(" {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")
	return out.String()
}

// ArrowFunction represents a lambda shorthand: x => x * 2
type ArrowFunction struct {
	Parameters []*ast.Identifier
	Body       ast.Expression // Single expression body
	Env        *Environment
}

func (af *ArrowFunction) Type() ObjectType { return ARROW_FUNCTION_OBJ }
func (af *ArrowFunction) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range af.Parameters {
		params = append(params, p.String())
	}
	if len(af.Parameters) == 1 {
		out.WriteString(af.Parameters[0].String())
	} else {
		out.WriteString("(")
		out.WriteString(strings.Join(params, ", "))
		out.WriteString(")")
	}
	out.WriteString(" => ")
	out.WriteString(af.Body.String())
	return out.String()
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

// Char represents a single character (rune value stored as int32)
type Char struct {
	Value rune
}

func (c *Char) Type() ObjectType { return CHAR_OBJ }
func (c *Char) Inspect() string  { return fmt.Sprintf("'%c'", c.Value) }
func (c *Char) HashKey() HashKey {
	return HashKey{Type: c.Type(), Value: uint64(c.Value)}
}

// Byte represents a single byte (0-255)
type Byte struct {
	Value byte
}

func (b *Byte) Type() ObjectType { return BYTE_OBJ }
func (b *Byte) Inspect() string  { return fmt.Sprintf("0x%02X", b.Value) }
func (b *Byte) HashKey() HashKey {
	return HashKey{Type: b.Type(), Value: uint64(b.Value)}
}

// Rune represents a Unicode code point (alias for Char but explicit)
type Rune struct {
	Value rune
}

func (r *Rune) Type() ObjectType { return RUNE_OBJ }
func (r *Rune) Inspect() string  { return fmt.Sprintf("'%c' (U+%04X)", r.Value, r.Value) }
func (r *Rune) HashKey() HashKey {
	return HashKey{Type: r.Type(), Value: uint64(r.Value)}
}

// Enum represents an enum type definition
type Enum struct {
	Name   string
	Values map[string]int64 // Maps enum variant names to their integer values
}

func (e *Enum) Type() ObjectType { return ENUM_OBJ }
func (e *Enum) Inspect() string {
	var out bytes.Buffer
	out.WriteString("enum ")
	out.WriteString(e.Name)
	out.WriteString(" { ")
	variants := []string{}
	for name, val := range e.Values {
		variants = append(variants, fmt.Sprintf("%s = %d", name, val))
	}
	out.WriteString(strings.Join(variants, ", "))
	out.WriteString(" }")
	return out.String()
}

// EnumValue represents a specific enum variant value
type EnumValue struct {
	EnumName  string // The enum type name
	ValueName string // The variant name
	Value     int64  // The integer value
}

func (ev *EnumValue) Type() ObjectType { return ENUM_VALUE_OBJ }
func (ev *EnumValue) Inspect() string  { return fmt.Sprintf("%s.%s", ev.EnumName, ev.ValueName) }
func (ev *EnumValue) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(ev.EnumName + "." + ev.ValueName))
	return HashKey{Type: ev.Type(), Value: h.Sum64()}
}

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

type Array struct {
	Elements []Object
}

func (ao *Array) Type() ObjectType { return ARRAY_OBJ }
func (ao *Array) Inspect() string {
	var out bytes.Buffer
	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out bytes.Buffer
	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

type Struct struct {
	Name   string
	Fields []string
}

func (s *Struct) Type() ObjectType { return STRUCT_OBJ }
func (s *Struct) Inspect() string {
	return "struct " + s.Name
}

type StructInstance struct {
	Struct *Struct
	Fields map[string]Object
}

func (si *StructInstance) Type() ObjectType { return INSTANCE_OBJ }
func (si *StructInstance) Inspect() string {
	var out bytes.Buffer
	out.WriteString(si.Struct.Name)
	out.WriteString(" { ")
	pairs := []string{}
	for k, v := range si.Fields {
		pairs = append(pairs, k+": "+v.Inspect())
	}
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString(" }")
	return out.String()
}

type Environment struct {
	store  map[string]Object
	consts map[string]bool // tracks which variables are const
	outer  *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	c := make(map[string]bool)
	return &Environment{store: s, consts: c, outer: nil}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

// SetConst sets a constant variable that cannot be reassigned
func (e *Environment) SetConst(name string, val Object) Object {
	e.store[name] = val
	e.consts[name] = true
	return val
}

// ToHash converts the environment's store to a Hash object
func (e *Environment) ToHash() *Hash {
	pairs := make(map[HashKey]HashPair)
	for name, val := range e.store {
		key := &String{Value: name}
		pairs[key.HashKey()] = HashPair{Key: key, Value: val}
	}
	return &Hash{Pairs: pairs}
}

// IsConst checks if a variable is a constant
func (e *Environment) IsConst(name string) bool {
	if isConst, ok := e.consts[name]; ok && isConst {
		return true
	}
	if e.outer != nil {
		return e.outer.IsConst(name)
	}
	return false
}

// Update updates an existing variable in the environment chain
func (e *Environment) Update(name string, val Object) (Object, bool) {
	_, ok := e.store[name]
	if ok {
		e.store[name] = val
		return val, true
	}
	if e.outer != nil {
		return e.outer.Update(name, val)
	}
	return nil, false
}

// Break represents a break statement
type Break struct{}

func (b *Break) Type() ObjectType { return BREAK_OBJ }
func (b *Break) Inspect() string  { return "break" }

// Continue represents a continue statement
type Continue struct{}

func (c *Continue) Type() ObjectType { return CONTINUE_OBJ }
func (c *Continue) Inspect() string  { return "continue" }

// Range represents a range like 0..10
type Range struct {
	Start int64
	End   int64
}

func (r *Range) Type() ObjectType { return RANGE_OBJ }
func (r *Range) Inspect() string  { return fmt.Sprintf("%d..%d", r.Start, r.End) }

// TypeChecker provides type validation utilities
// CheckType validates if an object matches the expected type annotation
func CheckType(obj Object, typeAnn *ast.TypeAnnotation) bool {
	if typeAnn == nil {
		return true // No type annotation means any type is allowed
	}

	typeName := typeAnn.TypeName

	// Handle 'any' type - matches anything
	if typeName == "any" {
		return true
	}

	// Handle array types
	if typeAnn.IsArray {
		arr, ok := obj.(*Array)
		if !ok {
			return false
		}
		// If element type is specified, check all elements
		if typeAnn.ElementType != nil {
			for _, elem := range arr.Elements {
				if !CheckType(elem, typeAnn.ElementType) {
					return false
				}
			}
		}
		return true
	}

	// Handle map types
	if typeAnn.KeyType != nil && typeAnn.ElementType != nil {
		_, ok := obj.(*Hash)
		return ok // For now, just check if it's a hash; deeper checking can be added
	}

	// Handle basic types
	switch typeName {
	case "int":
		_, ok := obj.(*Integer)
		return ok
	case "float":
		// Allow int as float for convenience (like Go)
		_, isFloat := obj.(*Float)
		_, isInt := obj.(*Integer)
		return isFloat || isInt
	case "string":
		_, ok := obj.(*String)
		return ok
	case "bool":
		_, ok := obj.(*Boolean)
		return ok
	case "char":
		// Char can be Char type or single-character string
		if _, ok := obj.(*Char); ok {
			return true
		}
		s, ok := obj.(*String)
		if ok && len(s.Value) == 1 {
			return true
		}
		return false
	case "byte":
		_, ok := obj.(*Byte)
		return ok
	case "rune":
		// Rune can be Rune type, Char type, or Integer
		if _, ok := obj.(*Rune); ok {
			return true
		}
		if _, ok := obj.(*Char); ok {
			return true
		}
		if _, ok := obj.(*Integer); ok {
			return true
		}
		return false
	case "array":
		_, ok := obj.(*Array)
		return ok
	case "map":
		_, ok := obj.(*Hash)
		return ok
	case "void":
		_, ok := obj.(*Null)
		return ok
	default:
		// Custom type (struct or enum) - check if it's a struct instance with matching name
		si, ok := obj.(*StructInstance)
		if ok && si.Struct.Name == typeName {
			return true
		}
		// Check for enum values
		ev, ok := obj.(*EnumValue)
		if ok && ev.EnumName == typeName {
			return true
		}
		return false
	}
}

// TypeName returns the type name of an object as a string
func TypeName(obj Object) string {
	switch obj := obj.(type) {
	case *Integer:
		return "int"
	case *Float:
		return "float"
	case *String:
		return "string"
	case *Char:
		return "char"
	case *Byte:
		return "byte"
	case *Rune:
		return "rune"
	case *Boolean:
		return "bool"
	case *Array:
		return "array"
	case *Hash:
		return "map"
	case *Null:
		return "void"
	case *Function:
		return "function"
	case *ArrowFunction:
		return "function"
	case *StructInstance:
		return obj.Struct.Name
	case *EnumValue:
		return obj.EnumName
	case *Enum:
		return "enum"
	default:
		return string(obj.Type())
	}
}
