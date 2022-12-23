package command_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kodefluence/altair/module/controller"
	"github.com/kodefluence/altair/plugin/oauth/entity"
	"github.com/kodefluence/altair/plugin/oauth/module/application/controller/command"
	"github.com/kodefluence/altair/plugin/oauth/module/application/controller/command/mock"
	"github.com/kodefluence/altair/testhelper"
	"github.com/kodefluence/altair/util"
	"github.com/kodefluence/monorepo/jsonapi"
	"github.com/kodefluence/monorepo/kontext"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestCreateApplication(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Given owner_id and scope, when command is executed then it should create the applications", func(t *testing.T) {
		cmd := &cobra.Command{
			Use: "test",
		}

		oauthApplicationJSON := entity.OauthApplicationJSON{
			OwnerID:     util.ValueToPointer(1),
			OwnerType:   util.ValueToPointer("confidential"),
			Description: util.ValueToPointer("test"),
			Scopes:      util.ValueToPointer("read write"),
		}

		applicationManager := mock.NewMockApplicationManager(mockCtrl)
		applicationManager.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(ktx kontext.Context, e entity.OauthApplicationJSON) (entity.OauthApplicationJSON, jsonapi.Errors) {
			assert.Equal(t, "test", util.PointerToValue(e.Description))
			assert.Equal(t, 1, util.PointerToValue(e.OwnerID))
			assert.Equal(t, "read write", util.PointerToValue(e.Scopes))
			assert.Equal(t, "confidential", util.PointerToValue(e.OwnerType))
			return oauthApplicationJSON, nil
		})

		command := command.NewCreateOauthApplication(applicationManager)

		appController := controller.Provide(nil, nil, cmd)
		appController.InjectCommand(command)

		// Given
		cmd.SetArgs([]string{"oauth/application:create", "--owner-id", "1", "--scope", "read write", "--owner-type", "confidential", "--desc", "test"})

		// When
		err := cmd.Execute()

		// Then
		assert.Nil(t, err)
	})

	t.Run("Given owner_id and scope, when there is error in command execution then it should print the error", func(t *testing.T) {
		cmd := &cobra.Command{
			Use: "test",
		}

		applicationManager := mock.NewMockApplicationManager(mockCtrl)
		applicationManager.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(ktx kontext.Context, e entity.OauthApplicationJSON) (entity.OauthApplicationJSON, jsonapi.Errors) {
			assert.Equal(t, "test", util.PointerToValue(e.Description))
			assert.Equal(t, 1, util.PointerToValue(e.OwnerID))
			assert.Equal(t, "read write", util.PointerToValue(e.Scopes))
			assert.Equal(t, "confidential", util.PointerToValue(e.OwnerType))
			return entity.OauthApplicationJSON{}, testhelper.ErrInternalServer()
		})

		command := command.NewCreateOauthApplication(applicationManager)

		appController := controller.Provide(nil, nil, cmd)
		appController.InjectCommand(command)

		// Given
		cmd.SetArgs([]string{"oauth/application:create", "--owner-id", "1", "--scope", "read write", "--owner-type", "confidential", "--desc", "test"})

		// When
		err := cmd.Execute()

		// Then
		assert.Nil(t, err)
	})
}
