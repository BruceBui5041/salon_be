package modelhelper

import (
	"fmt"
	"reflect"
	"sync"
)

var (
	modelRegistry    = make(map[string]reflect.Type)
	responseRegistry = make(map[string]reflect.Type)
	registryMutex    sync.RWMutex
)

// RegisterModel adds a model to the registry for automatic type lookup
func RegisterModel(model interface{}) {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if m, ok := reflect.New(t).Interface().(interface{ TableName() string }); ok {
		registryMutex.Lock()
		modelRegistry[m.TableName()] = t
		registryMutex.Unlock()
	}
}

// RegisterResponseType associates a response type with a table name in the registry
func RegisterResponseType(tableName string, responseType interface{}) {
	t := reflect.TypeOf(responseType)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	registryMutex.Lock()
	responseRegistry[tableName] = t
	registryMutex.Unlock()
}

// GetModelType retrieves the registered model type for a given table name
func GetModelType(modelName string) (reflect.Type, error) {
	registryMutex.RLock()
	t, ok := modelRegistry[modelName]
	registryMutex.RUnlock()

	if !ok {
		return nil, fmt.Errorf("model not being registered: %s", modelName)
	}
	return t, nil
}

// GetSearchModelResponseType retrieves the registered response type for a given table name
func GetSearchModelResponseType(modelName string) (reflect.Type, error) {
	registryMutex.RLock()
	t, ok := responseRegistry[modelName]
	registryMutex.RUnlock()

	if !ok {
		return GetModelType(modelName)
	}
	return t, nil
}
