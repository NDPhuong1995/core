package beans

import (
	"fmt"
	"reflect"
	"sync"
)

var mutex sync.Mutex
var beansTypeRegistries = make(map[string]reflect.Type)
var beansInstanceRegistries = make(map[string]interface{})

type ScopeBean int8

const (
	ScopeSingleton = 1
	ScopePrototype = 2
)

func RegistryBean(beanName string, value interface{}) error {
	if _, err := getBeanType(beanName); err == nil {
		return fmt.Errorf("duplicate bean with name %s", beanName)
	}
	beansTypeRegistries[beanName] = reflect.TypeOf(value)
	return nil
}

func GetBean(beanName string, scope ScopeBean) (interface{}, error) {
	beanType, err := getBeanType(beanName)
	if scope == ScopePrototype {
		return reflect.New(beanType).Elem().Interface(), nil
	}
	mutex.Lock()
	defer mutex.Unlock()
	if bean, isSuccess := beansInstanceRegistries[beanName]; isSuccess {
		return bean, nil
	}
	if err != nil  {
		return nil, err
	}
	beansInstanceRegistries[beanName] = reflect.New(beanType).Elem().Interface()
	return beansInstanceRegistries[beanName], nil
}

func InvokeBeanMethod(bean interface{}, methodName string, args ...interface{}) ([]interface{}, error) {
	methodValue := reflect.ValueOf(bean).MethodByName(methodName)
	if methodValue.IsZero() {
		return nil, fmt.Errorf("no method %s in bean", methodName)
	}
	methodType := methodValue.Type()
	numIn := methodType.NumIn()
	if numIn > len(args) {
		return nil, fmt.Errorf("method %s must have minimum %d params have %d", methodName, numIn, len(args))
	}
	if numIn != len(args) && !methodType.IsVariadic(){
		return nil, fmt.Errorf("method %s must have %d params have %d", methodName, numIn, len(args))
	}
	params := make([]reflect.Value, len(args))
	for i := 0; i < len(args); i++ {
		var inType reflect.Type
		if methodType.IsVariadic() && i >= numIn-1 {
			inType = methodType.In(numIn - 1).Elem()
		} else {
			inType = methodType.In(i)
		}
		argValue := reflect.ValueOf(args[i])
		if !argValue.IsValid() {
			return nil, fmt.Errorf("method %s. param[%d] must be %s have %s", methodName, i, inType, argValue.String())
		}
		argType := argValue.Type()
		if argType.ConvertibleTo(inType) {
			params[i] = argValue.Convert(inType)
		} else {
			return nil, fmt.Errorf("method %s. param[%d] must be %s have %s", methodName, i, inType, argType)
		}
	}
	values := methodValue.Call(params)
	var results = make([]interface{}, len(values))
	for i := 0; i < len(results); i++ {
		results[i] = values[i].Interface()
	}
	return results, nil
}

func getBeanType(beanName string) (reflect.Type, error) {
	if bean, isSuccess := beansTypeRegistries[beanName]; isSuccess {
		return bean, nil
	}
	return nil, fmt.Errorf("no bean with name %s", beanName)
}