package jsonapi

import "github.com/kodefluence/monorepo/exception"

type Option func(*Body)

func WithData(data interface{}) Option {
	return func(b *Body) {
		b.Data = data
	}
}

func WithErrors(errors Errors) Option {
	return func(b *Body) {
		if errors != nil {
			b.Errors = errors
		}
	}
}

type ExceptionOption func(*Error)

func WithException(code string, status int, exc exception.Exception, options ...ExceptionOption) Option {
	return WithExceptionMeta(code, status, exc, nil, options...)
}

func WithExceptionMeta(code string, status int, exc exception.Exception, meta Meta, options ...ExceptionOption) Option {
	return func(b *Body) {
		err := Error{
			Title:  exc.Title(),
			Detail: exc.Detail(),
			Code:   code,
			Status: status,
			Meta:   meta,
		}
		for _, opt := range options {
			opt(&err)
		}
		b.Errors = append(b.Errors, err)
	}
}

// WithSourcePointer creates a Source with pointer field set
func WithSourcePointer(pointer string) ExceptionOption {
	return func(err *Error) {
		if err.Source == nil {
			err.Source = &Source{}
		}
		err.Source.Pointer = pointer
	}
}

// WithSourceParameter creates a Source with parameter field set
func WithSourceParameter(parameter string) ExceptionOption {
	return func(err *Error) {
		if err.Source == nil {
			err.Source = &Source{}
		}
		err.Source.Parameter = parameter
	}
}

// WithSourceHeader creates a Source with header field set
func WithSourceHeader(header string) ExceptionOption {
	return func(err *Error) {
		if err.Source == nil {
			err.Source = &Source{}
		}
		err.Source.Header = header
	}
}

func WithMeta(key string, field interface{}) Option {
	return func(b *Body) {
		if b.Meta == nil {
			b.Meta = Meta{}
		}
		b.Meta[key] = field
	}
}
