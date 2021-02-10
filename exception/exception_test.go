package exception_test

import (
	"errors"
	"testing"

	"github.com/codefluence-x/monorepo/exception"
	"github.com/stretchr/testify/assert"
)

func TestException(t *testing.T) {

	t.Run("does return an exception", func(t *testing.T) {
		err := errors.New("unexpected error")
		exceptionType := exception.NotFound
		detail := "data not found in the databases because of deletion"
		title := "data is not exists"

		exc := exception.Throw(err, exception.WithType(exceptionType), exception.WithTitle(title), exception.WithDetail(detail))
		assert.Equal(t, exceptionType, exc.Type())
		assert.Equal(t, exceptionType.String(), exc.Type().String())
		assert.Equal(t, detail, exc.Detail())
		assert.Equal(t, title, exc.Title())
		assert.Equal(t, err.Error(), exc.Error())
	})

	t.Run("Type", func(t *testing.T) {
		assert.Equal(t, "unexpected", exception.Unexpected.String())
		assert.Equal(t, "not found", exception.NotFound.String())
		assert.Equal(t, "duplicated", exception.Duplicated.String())
		assert.Equal(t, "bad input", exception.BadInput.String())
	})
}
