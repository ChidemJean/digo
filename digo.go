package digo

import (
	"fmt"
	"reflect"
)

var container *Container

type Scope int

const (
	Singleton Scope = iota
	Transient
)

type provider struct {
	constructor reflect.Value
	scope       Scope
	instance    reflect.Value
	outType     reflect.Type
}

type Container struct {
	providers map[reflect.Type]*provider
	aliases   map[reflect.Type]reflect.Type
}

func New() *Container {
	c := &Container{
		providers: make(map[reflect.Type]*provider),
		aliases:   make(map[reflect.Type]reflect.Type),
	}
	container = c
	return c
}

func (c *Container) Register(constructor interface{}, scope Scope) {
	val := reflect.ValueOf(constructor)
	typ := val.Type()

	if typ.Kind() != reflect.Func || typ.NumOut() != 1 {
		fmt.Printf("Constructor must be a function with exactly one return")
	}

	outType := typ.Out(0)
	c.providers[outType] = &provider{
		constructor: val,
		scope:       scope,
		outType:     outType,
	}
}

func (c *Container) RegisterInterface(interfacePtr interface{}, implementation interface{}) error {
	ifaceType := reflect.TypeOf(interfacePtr).Elem()
	implType := reflect.TypeOf(implementation)

	if !implType.Implements(ifaceType) {
		return fmt.Errorf("%v does not implement %v", implType, ifaceType)
	}
	c.aliases[ifaceType] = implType

	return nil
}

func Resolve[T any]() (T, error) {
	t := reflect.TypeOf((*T)(nil)).Elem()
	val, err := container.resolve(t)
	if !val.IsValid() || err != nil {
		var zero T
		return zero, err
	}
	return val.Interface().(T), nil
}

func (c *Container) resolve(t reflect.Type) (reflect.Value, error) {
	if t.Kind() == reflect.Interface {
		if impl, ok := c.aliases[t]; ok {
			t = impl
		} else {
			return reflect.Value{}, fmt.Errorf("no implementation registered for interface: %v", t)
		}
	}

	prov, ok := c.providers[t]
	if !ok {
		return reflect.Value{}, fmt.Errorf("no provider found for type: %v", t)
	}

	if prov.scope == Singleton && prov.instance.IsValid() {
		return prov.instance, nil
	}

	args := []reflect.Value{}
	for i := 0; i < prov.constructor.Type().NumIn(); i++ {
		paramType := prov.constructor.Type().In(i)
		arg, err := c.resolve(paramType)
		if err != nil {
			arg = reflect.Value{}
		}
		args = append(args, arg)
	}

	result := prov.constructor.Call(args)[0]

	if prov.scope == Singleton {
		prov.instance = result
	}

	return result, nil
}
