package monorepo

import "context"

// Context in monorepo wrap default context implementation with Get and Set which can carry a value from one flow
type Context interface {
	Ctx() context.Context

	// Get retrieves data from the context.
	Get(key string) interface{}

	// Set saves data in the context.
	Set(key string, val interface{})
}
