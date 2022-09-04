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

func WithException(code string, status int, exc exception.Exception) Option {
	return WithExceptionMeta(code, status, exc, nil)
}

func WithExceptionMeta(code string, status int, exc exception.Exception, meta Meta) Option {
	return func(b *Body) {
		err := Error{
			Title:  exc.Title(),
			Detail: exc.Detail(),
			Code:   code,
			Status: status,
			Meta:   meta,
		}
		b.Errors = append(b.Errors, err)
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
