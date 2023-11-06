package object

func NewEnvironment(e *Environment) *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: e}
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

func (e *Environment) Reset(name string, val Object) (Object, bool) {
	var ok bool
	_, ok = e.store[name]
	if ok {
		e.store[name] = val
	}

	if !ok && e.outer != nil {
		_, ok = e.outer.Reset(name, val)
	}

	if !ok {
		e.store[name] = val
		ok = true
	}
	return val, ok
}
