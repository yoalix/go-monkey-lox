package evaluator

import (
	"fmt"
	"go-compiler/main/ast"
	"go-compiler/main/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node, env *object.Environment) object.Object {

	switch node := node.(type) {
	case *ast.Program:
		return evalProgramStatements(node, env)
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ast.AssignStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Reset(node.Name.Value, val)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.NumberLiteral:
		return &object.Number{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right, node.Token.Line)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right, node.Token.Line)
	case *ast.BlockStatment:
		return evalBlockStatements(node.Statements, env)
	case *ast.IfExpression:
		condition := Eval(node.Condition, env)
		if isError(condition) {
			return condition
		}
		if condition != FALSE && condition != NULL {
			return Eval(node.Then, env)
		} else if node.Else != nil {
			return Eval(node.Else, env)
		} else {
			return NULL
		}
	case *ast.WhileExpression:
		for {

			condition := Eval(node.Condition, env)
			if isError(condition) {
				return condition
			}
			if !isTruthy(condition) {
				break
			}
			eval := Eval(node.Body, env)
			if isError(eval) {
				return eval
			}
		}
		return NULL
	case *ast.ReturnStatement:
		ret := Eval(node.ReturnValue, env)
		if isError(ret) {
			return ret
		}
		return &object.ReturnValue{Value: ret}
	case *ast.FunctionLiteral:
		return &object.Function{
			Parameters: node.Parameters,
			Body:       node.Body,
			Env:        env,
		}
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env, node.Token.Line)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return evalCall(function, args, node.Token.Line)
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env, node.Token.Line)
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

		return evalIndexExpression(left, index, node.Token.Line)
	case *ast.HashLiteral:

		return evalHashLiteral(node, env)
	}

	return nil
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if ident, ok := env.Get(node.Value); ok {
		return ident
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	return newError("[line %v] identifier not found: %v", node.Token.Line, node.Value)
}

func evalProgramStatements(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object
	for _, stmt := range program.Statements {
		result = Eval(stmt, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue.Value
		}
	}
	return result
}
func evalBlockStatements(statements []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range statements {
		result = Eval(stmt, env)
		if result != nil {
			if result.Type() == object.RETURN_OBJ || result.Type() == object.ERROR_OBJ {
				return result
			}
		}
	}
	return result
}

// func evalWhileExpression(condition, body object.Object, env *object.Environment)

func evalExpressions(expressions []ast.Expression, env *object.Environment, line int) []object.Object {
	var result []object.Object
	for _, expr := range expressions {
		evaluated := Eval(expr, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func evalCall(fn object.Object, args []object.Object, line int) object.Object {
	switch fn := fn.(type) {
	case *object.Function:

		extendedEnv := extendFunctionEnv(fn, args)
		eval := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(eval)
	case *object.Builtin:
		evaluated := fn.Fn(args...)
		if evaluated, ok := evaluated.(*object.Error); ok {
			evaluated.Message = fmt.Sprintf("[line %v] %v", line, evaluated.Message)
		}
		return evaluated
	default:
		return newError("[line %v] not a function: %s", line, fn.Type())
	}

}

func extendFunctionEnv(
	fn *object.Function,
	args []object.Object,
) *object.Environment {
	env := object.NewEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func evalIndexExpression(left, index object.Object, line int) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.NUMBER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:

		return evalHashIndexExpression(left, index, line)
	default:
		return newError("[line %v] index operator not supported: %s", line, left.Type())
	}

}

func evalArrayIndexExpression(left, index object.Object) object.Object {
	arr := left.(*object.Array).Elements
	idx := index.(*object.Number).Value
	max := float64(len(arr) - 1)
	if idx < 0 || idx > max {
		return NULL
	}
	return arr[int(idx)]
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	hash := &object.Hash{Pairs: make(map[object.HashKey]object.HashPair)}
	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}
		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("[line %v] unusable as hash key: %s", node.Token.Line, key.Type())
		}
		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}
		hash.Pairs[hashKey.HashKey()] = object.HashPair{Key: key, Value: value}
	}
	return hash
}
func evalHashIndexExpression(left, index object.Object, line int) object.Object {
	hash := left.(*object.Hash)
	hashKey, ok := index.(object.Hashable)
	if !ok {
		return newError("[line %v] unusable as hash key: %s", line, index.Type())
	}

	val, ok := hash.Pairs[hashKey.HashKey()]
	if !ok {
		return NULL
	}
	return val.Value
}

func evalPrefixExpression(op string, right object.Object, line int) object.Object {
	switch op {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusOperatorExpression(right, line)
	default:
		return newError("[line %v] Unknown operator: %s%s", line, op, right.Type())
	}
}

func evalInfixExpression(op string, left, right object.Object, line int) object.Object {
	// case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
	// 	return evalBooleanInfixExpression(op, left, right)
	switch {
	case left.Type() == object.NUMBER_OBJ && right.Type() == object.NUMBER_OBJ:
		return evalNumberInfixExpression(op, left, right, line)
	case (left.Type() == object.STRING_OBJ || left.Type() == object.NUMBER_OBJ) &&
		(right.Type() == object.NUMBER_OBJ || right.Type() == object.STRING_OBJ):
		return evalStringInfixExpression(op, left, right, line)
	case left.Type() != right.Type():
		return newError("[line %v] type mismatch: %s %s %s",
			line, left.Type(), op, right.Type())
	case op == "==":
		return nativeBoolToBooleanObject(left == right)
	case op == "!=":
		return nativeBoolToBooleanObject(left != right)
	case op == "and":
		return nativeBoolToBooleanObject(left.(*object.Boolean).Value && right.(*object.Boolean).Value)
	case op == "or":
		return nativeBoolToBooleanObject(left.(*object.Boolean).Value || right.(*object.Boolean).Value)
	default:
		return newError("[line %v] unknown operator: %s %s %s", line, left.Type(), op, right.Type())
	}
}

func evalNumberInfixExpression(op string, left, right object.Object, line int) object.Object {
	leftVal := left.(*object.Number).Value
	rightVal := right.(*object.Number).Value
	switch op {
	case "+":
		return &object.Number{Value: leftVal + rightVal}
	case "-":
		return &object.Number{Value: leftVal - rightVal}
	case "*":
		return &object.Number{Value: leftVal * rightVal}
	case "/":
		return &object.Number{Value: leftVal / rightVal}
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("[line %v] unknown operator: %s %s %s",
			line, left.Type(), op, right.Type())
	}
}

func evalStringInfixExpression(op string, left, right object.Object, line int) object.Object {
	if op == "+" {
		str := fmt.Sprintf("%v%v", left.Inspect(), right.Inspect())
		return &object.String{Value: str}
	}
	leftVal := left.Inspect()
	rightVal := right.Inspect()
	switch op {
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	}
	return newError("[line %v] unknown operator: %s %s %s",
		line, left.Type(), op, right.Type())
}

// func evalBooleanInfixExpression(op string, left, right object.Object) object.Object {
// 	leftVal := left.(*object.Boolean).Value
// 	rightVal := right.(*object.Boolean).Value
// 	switch op {
// 	case "==":
// 		return &object.Boolean{Value: leftVal == rightVal}
// 	case "!=":
// 		return &object.Boolean{Value: leftVal != rightVal}
// 	default:
// 		return NULL
// 	}
// }

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return NULL
	default:
		return FALSE
	}
}

func evalMinusOperatorExpression(right object.Object, line int) object.Object {
	switch right.Type() {
	case object.NUMBER_OBJ:
		return &object.Number{Value: -right.(*object.Number).Value}
	default:
		return newError("[line %v] unknown operator: -%s", line, right.Type())
	}

}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func isTruthy(val interface{}) bool {
	if val == nil {
		return false
	} else if b, ok := val.(bool); ok {
		return b
	}
	return true
}
