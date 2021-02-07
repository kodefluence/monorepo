package kontext

import (
	"context"
	"sync"
)

// Context in wrap default context implementation with Get and Set which can carry a value from one flow
type Context interface {
	Ctx() context.Context
	Get(key string) (interface{}, bool)
	GetWithoutCheck(key string) interface{}
	Set(key string, val interface{})
}

// Kontext implement kontext.Context
type Kontext struct {
	ctx    context.Context
	values *sync.Map
}

// Fabricate Context
func Fabricate(opts ...Option) *Kontext {
	var config Config
	config.ctx = context.Background()

	for _, opt := range opts {
		opt(&config)
	}

	return &Kontext{
		ctx:    config.ctx,
		values: &sync.Map{},
	}
}

// Ctx carried default built in context.Context
func (k *Kontext) Ctx() context.Context {
	return k.ctx
}

// Get data from the context. It will return value and a boolean to indicate if the value is exists or not
func (k *Kontext) Get(key string) (interface{}, bool) {
	return k.values.Load(key)
}

// GetWithoutCheck will return value from the Context without returning boolean indicator if the value is exists or not
func (k *Kontext) GetWithoutCheck(key string) interface{} {
	v, _ := k.Get(key)
	return v
}

// Set data in the context based on it's key value
func (k *Kontext) Set(key string, val interface{}) {
	k.values.Store(key, val)
}
