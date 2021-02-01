package entity_test

import (
	"testing"
	"time"

	"github.com/codefluence-x/altair/provider/plugin/oauth/entity"
	"github.com/stretchr/testify/assert"
)

func TestEntity(t *testing.T) {
	t.Run("ErrorObject", func(t *testing.T) {
		errorObject := entity.ErrorObject{
			Code:    "ERR001",
			Message: "Some error messages",
		}

		assert.Equal(t, "Error: Some error messages, Code: ERR001", errorObject.Error())
	})

	t.Run("Error", func(t *testing.T) {
		err := entity.Error{
			Errors: []entity.ErrorObject{
				entity.ErrorObject{
					Code:    "ERR001",
					Message: "Some error messages",
				},
				entity.ErrorObject{
					Code:    "ERR002",
					Message: "Some error messages",
				},
			},
		}

		assert.Equal(t, "Service error because of:\nError: Some error messages, Code: ERR001\nError: Some error messages, Code: ERR002\n", err.Error())
	})

}

func TestOauthPlugin(t *testing.T) {

	oauthPlugin := entity.OauthPlugin{}
	oauthPlugin.Config.Database = "main_database"
	oauthPlugin.Config.AccessTokenTimeoutRaw = "24h"
	oauthPlugin.Config.AuthorizationCodeTimeoutRaw = "24h"
	oauthPlugin.Config.RefreshToken.Timeout = "24h"

	t.Run("DatabaseInstance", func(t *testing.T) {
		t.Run("Return database instance", func(t *testing.T) {
			assert.Equal(t, "main_database", oauthPlugin.DatabaseInstance())
		})
	})

	t.Run("AccessTokenTimeout", func(t *testing.T) {
		t.Run("Right format", func(t *testing.T) {
			t.Run("Return duration", func(t *testing.T) {
				duration, err := oauthPlugin.AccessTokenTimeout()
				assert.Nil(t, err)
				assert.Equal(t, time.Hour*24, duration)
			})
		})

		t.Run("Wrong format", func(t *testing.T) {
			t.Run("Return error", func(t *testing.T) {
				oauthPlugin := entity.OauthPlugin{}
				oauthPlugin.Config.Database = "main_database"
				oauthPlugin.Config.AccessTokenTimeoutRaw = "abc"
				oauthPlugin.Config.AuthorizationCodeTimeoutRaw = "24h"

				_, err := oauthPlugin.AccessTokenTimeout()
				assert.NotNil(t, err)
			})
		})
	})

	t.Run("AuthorizationCodeTimeout", func(t *testing.T) {
		t.Run("Right format", func(t *testing.T) {
			t.Run("Return duration", func(t *testing.T) {
				duration, err := oauthPlugin.AuthorizationCodeTimeout()
				assert.Nil(t, err)
				assert.Equal(t, time.Hour*24, duration)
			})
		})

		t.Run("Wrong format", func(t *testing.T) {
			t.Run("Return error", func(t *testing.T) {
				oauthPlugin := entity.OauthPlugin{}
				oauthPlugin.Config.Database = "main_database"
				oauthPlugin.Config.AccessTokenTimeoutRaw = "24h"
				oauthPlugin.Config.AuthorizationCodeTimeoutRaw = "abc"

				_, err := oauthPlugin.AuthorizationCodeTimeout()
				assert.NotNil(t, err)
			})
		})
	})

	t.Run("RefreshTokenTimeout", func(t *testing.T) {
		t.Run("Right format", func(t *testing.T) {
			t.Run("Return duration", func(t *testing.T) {
				duration, err := oauthPlugin.RefreshTokenTimeout()
				assert.Nil(t, err)
				assert.Equal(t, time.Hour*24, duration)
			})
		})

		t.Run("Wrong format", func(t *testing.T) {
			t.Run("Return error", func(t *testing.T) {
				oauthPlugin := entity.OauthPlugin{}
				oauthPlugin.Config.RefreshToken.Timeout = "abc"

				_, err := oauthPlugin.RefreshTokenTimeout()
				assert.NotNil(t, err)
			})
		})
	})
}
