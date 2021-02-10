package kontext_test

import (
	"context"
	"testing"

	"github.com/codefluence-x/monorepo/kontext"
	"github.com/stretchr/testify/assert"
)

func TestKontext(t *testing.T) {

	t.Run("without option", func(t *testing.T) {
		ktx := kontext.Fabricate()
		key := "some-value"
		val := 100
		ktx.Set(key, val)

		returnedValue, exists := ktx.Get(key)
		assert.True(t, exists)
		assert.Equal(t, val, returnedValue)

		returnedValueWithoutCheck := ktx.GetWithoutCheck(key)
		assert.Equal(t, val, returnedValueWithoutCheck)

		_, exists = ktx.Get("non-exists-key")
		assert.False(t, exists)
	})

	t.Run("with default context option", func(t *testing.T) {
		ctx := context.Background()
		ktx := kontext.Fabricate(kontext.WithDefaultContext(ctx))

		assert.Equal(t, ctx, ktx.Ctx())
	})
}
