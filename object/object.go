package object

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"go-compiler/main/ast"
	"hash/fnv"
	"strings"
)

type ObjectType string

const (
	NUMBER_OBJ   = "NUMBER"
	STRING_OBJ   = "STRING"
	BOOLEAN_OBJ  = "BOOLEAN"
	NULL_OBJ     = "NULL"
	RETURN_OBJ   = "RETURN"
	FUNCITON_OBJ = "FUNCTION"
	BUILTIN_OBJ  = "BUILTIN"
	ARRAY_OBJ    = "ARRAY"
	HASH_OBJ     = "HASH"
	ERROR_OBJ    = "ERROR"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Number struct {
	Value float64
}

func (n *Number) Inspect() string  { return fmt.Sprintf("%v", n.Value) }
func (n *Number) Type() ObjectType { return NUMBER_OBJ }

type String struct {
	Value string
}

func (n *String) Inspect() string  { return fmt.Sprintf("%v", n.Value) }
func (n *String) Type() ObjectType { return STRING_OBJ }

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

type Null struct{}

func (n *Null) Inspect() string  { return "null" }
func (n *Null) Type() ObjectType { return NULL_OBJ }

type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) Inspect() string  { return r.Value.Inspect() }
func (r *ReturnValue) Type() ObjectType { return RETURN_OBJ }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatment
	Env        *Environment
}

func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, param := range f.Parameters {
		params = append(params, param.String())
	}

	out.WriteString("fn(" + strings.Join(params, ", ") + ") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")
	return out.String()
}
func (f *Function) Type() ObjectType { return FUNCITON_OBJ }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

func NewEnvironment(e *Environment) *Environment {
	s := make(map[string]Object)
	return &Environment{s, e}
}

type Environment struct {
	store map[string]Object
	outer *Environment
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, obj Object) Object {
	e.store[name] = obj
	return obj
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

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	var out bytes.Buffer
	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (a *Hash) Inspect() string {
	var out bytes.Buffer
	elements := []string{}
	for _, pairs := range a.Pairs {
		elements = append(elements, pairs.Key.Inspect()+":"+pairs.Value.Inspect())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("}")
	return out.String()
}

type HashKey struct {
	Type  ObjectType
	Value uint32
}

func (b *Boolean) HashKey() HashKey {
	var val uint32
	if b.Value {
		val = 1
	} else {
		val = 0
	}

	return HashKey{Type: b.Type(), Value: val}
}

func (n *Number) HashKey() HashKey {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint32(b, uint32(n.Value))
	h := fnv.New32a()
	h.Write(b)
	return HashKey{Type: n.Type(), Value: h.Sum32()}
	// return HashKey{Type: n.Type(), Value: n.Value}
}

func (s *String) HashKey() HashKey {
	h := fnv.New32a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum32()}
}

type Hashable interface {
	HashKey() HashKey
}
