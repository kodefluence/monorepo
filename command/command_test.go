package command_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kodefluence/monorepo/command"
	"github.com/kodefluence/monorepo/command/mock"

	"github.com/stretchr/testify/assert"
)

func TestCommand(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("monorepo", func(t *testing.T) {
		cmd := command.Fabricate()
		cmd.SetArgs([]string{"monorepo"})
		assert.Nil(t, cmd.Execute())
	})

	t.Run("custom", func(t *testing.T) {
		cmd := command.Fabricate(command.Config{
			Name:  "altair",
			Short: "Open Source API-Gateway",
		})
		cmd.SetArgs([]string{"altair"})
		assert.Nil(t, cmd.Execute())
	})

	t.Run("others", func(t *testing.T) {
		cmd := command.Fabricate()
		cmd.SetArgs([]string{"others"})

		scaffolding := mock.NewMockScaffold(mockCtrl)

		scaffolding.EXPECT().Use().Return("others")
		scaffolding.EXPECT().Short().Return("Others command")
		scaffolding.EXPECT().Example().Return("others [command]")
		scaffolding.EXPECT().Run([]string{})
		cmd.InjectCommand(scaffolding)
		assert.Nil(t, cmd.Execute())
	})
}
