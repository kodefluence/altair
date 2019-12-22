package entity_test

import (
	"testing"

	"github.com/codefluence-x/altair/entity"
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
