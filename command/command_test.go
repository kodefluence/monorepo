package command_test

import (
	"testing"

	"github.com/codefluence-x/monorepo/command"
	"github.com/codefluence-x/monorepo/monomock"
	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/assert"
)

func TestCommand(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("coronator", func(t *testing.T) {
		cmd := command.Fabricate()
		cmd.SetArgs([]string{"coronator"})
		assert.Nil(t, cmd.Execute())
	})

	t.Run("others", func(t *testing.T) {
		cmd := command.Fabricate()
		cmd.SetArgs([]string{"others"})

		scaffolding := monomock.NewMockCommandScaffold(mockCtrl)

		scaffolding.EXPECT().Use().Return("others")
		scaffolding.EXPECT().Short().Return("Others command")
		scaffolding.EXPECT().Example().Return("others [command]")
		scaffolding.EXPECT().Run([]string{})
		cmd.InjectCommand(scaffolding)
		assert.Nil(t, cmd.Execute())
	})
}
