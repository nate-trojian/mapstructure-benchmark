package mapstructurebenchmark

import (
	"encoding/json"
	"errors"
	"reflect"
)

type Registry struct {
	intType reflect.Type
	keyFunc func([]byte) (string, error)
	mapping map[string]reflect.Type
}

func (r *Registry) Register(key string, t reflect.Type) error {
	if !t.Implements(r.intType) {
		return errors.New("does not implement type")
	}
	r.mapping[key] = t
	return nil
}

func (r *Registry) FromJSON(data []byte) (interface{}, error) {
	key, err := r.keyFunc(data)
	if err != nil {
		return nil, err
	}
	t := r.mapping[key]
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	ret := reflect.New(t).Interface()
	if err := json.Unmarshal(data, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func NewRegistry(t reflect.Type, keyFunc func([]byte) (string, error)) *Registry {
	if t.Kind() != reflect.Interface {
		return nil
	}
	return &Registry{
		intType: t,
		keyFunc: keyFunc,
		mapping: make(map[string]reflect.Type),
	}
}
