package entity_test

import (
	"testing"

	"github.com/codefluence-x/altair/entity"
	"github.com/stretchr/testify/assert"
)

func TestEntity(t *testing.T) {
	t.Run("ResponseError", func(t *testing.T) {
		responseError := entity.ResponseError{
			Code:    "ERR001",
			Message: "Some error messages",
		}

		assert.Equal(t, "Error: Some error messages, Code: ERR001", responseError.Error())
	})

	t.Run("Error", func(t *testing.T) {
		responseError := entity.Error{
			Errors: []entity.ResponseError{
				entity.ResponseError{
					Code:    "ERR001",
					Message: "Some error messages",
				},
				entity.ResponseError{
					Code:    "ERR002",
					Message: "Some error messages",
				},
			},
		}

		assert.Equal(t, "Service error because of:\nError: Some error messages, Code: ERR001\nError: Some error messages, Code: ERR002\n", responseError.Error())
	})

}
