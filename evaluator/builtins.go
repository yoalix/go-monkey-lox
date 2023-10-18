package evaluator

import (
	"fmt"
	"go-compiler/main/object"
)

var builtins = map[string]*object.Builtin{
	"puts": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Printf("%v ", arg.Inspect())
			}
			fmt.Println()
			return NULL
		},
	},
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Number{Value: float64(len(arg.Elements))}
			case *object.String:
				return &object.Number{Value: float64(len(arg.Value))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) > 0 {
					return arg.Elements[0]
				}
				return NULL
			case *object.String:
				if len(arg.Value) > 0 {
					return &object.String{Value: string(arg.Value[0])}
				}
				return NULL
			default:
				return newError("argument to `first` not supported, got %s", args[0].Type())
			}
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) > 0 {
					return arg.Elements[len(arg.Elements)-1]
				}
				return NULL
			case *object.String:
				if len(arg.Value) > 0 {
					return &object.String{Value: string(arg.Value[len(arg.Value)-1])}
				}
				return NULL
			default:
				return newError("argument to `last` not supported, got %s", args[0].Type())
			}
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				length := len(arg.Elements)
				newElems := make([]object.Object, length)
				copy(newElems, arg.Elements)
				newElems = append(newElems, args[1])
				return &object.Array{Elements: newElems}
			default:
				return newError("argument to `push` must be ARRAY, got %s", args[0].Type())
			}
		},
	},
	"rest": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:

				length := len(arg.Elements)
				if length > 0 {

					newElems := make([]object.Object, length-1)
					copy(newElems, arg.Elements[1:length])
					return &object.Array{Elements: newElems}
				}
				return NULL
			default:
				return newError("argument to `rest` must be ARRAY, got %s", args[0].Type())
			}
		},
	},
}
