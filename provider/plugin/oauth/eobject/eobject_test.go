package eobject_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/codefluence-x/altair/provider/plugin/oauth/entity"
	"github.com/codefluence-x/altair/provider/plugin/oauth/eobject"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestErrorObject(t *testing.T) {
	t.Run("Wrap", func(t *testing.T) {
		t.Run("Wrap 1 or more error", func(t *testing.T) {
			ctx := context.WithValue(context.Background(), "track_id", uuid.New())
			err := eobject.Wrap(eobject.InternalServerError(ctx))

			assert.Equal(t, 1, len(err))
		})
	})

	t.Run("Internal server error", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "track_id", uuid.New())
		errorObject := eobject.InternalServerError(ctx)
		expectedErrorObject := entity.ErrorObject{
			Code:    "ERR0500",
			Message: fmt.Sprintf("Something is not right, help us fix this problem. Contribute to https://github.com/codefluence-x/altair. Or help us by give this code '%v' to the admin of this site.", ctx.Value("track_id")),
		}

		assert.Equal(t, expectedErrorObject.Code, errorObject.Code)
		assert.Equal(t, expectedErrorObject.Message, errorObject.Message)
		assert.Equal(t, expectedErrorObject.Error(), errorObject.Error())
	})

	t.Run("Bad request error", func(t *testing.T) {
		errorObject := eobject.BadRequestError("query parameter")
		expectedErrorObject := entity.ErrorObject{
			Code:    "ERR0400",
			Message: fmt.Sprintf("You've send malformed request in your `%s`", "query parameter"),
		}

		assert.Equal(t, expectedErrorObject.Code, errorObject.Code)
		assert.Equal(t, expectedErrorObject.Message, errorObject.Message)
		assert.Equal(t, expectedErrorObject.Error(), errorObject.Error())
	})

	t.Run("Not found error", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, "track_id", "1234567890")
		errorObject := eobject.NotFoundError(ctx, "some entity")
		expectedErrorObject := entity.ErrorObject{
			Code:    "ERR0404",
			Message: fmt.Sprintf("Resource of `%s` is not found, please report to admin of this site with this code `%v` if you think this is an error.", "some entity", ctx.Value("track_id")),
		}

		assert.Equal(t, expectedErrorObject.Code, errorObject.Code)
		assert.Equal(t, expectedErrorObject.Message, errorObject.Message)
		assert.Equal(t, expectedErrorObject.Error(), errorObject.Error())
	})

	t.Run("Forbidden error", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, "track_id", "1234567890")
		errorObject := eobject.ForbiddenError(ctx, "some entity", "your requested scope not exists in this applcation")
		expectedErrorObject := entity.ErrorObject{
			Code:    "ERR0403",
			Message: fmt.Sprintf("Resource of `%s` is forbidden to be accessed, because of: %s. Please report to admin of this site with this code `%v` if you think this is an error.", "some entity", "your requested scope not exists in this applcation", ctx.Value("track_id")),
		}

		assert.Equal(t, expectedErrorObject.Code, errorObject.Code)
		assert.Equal(t, expectedErrorObject.Message, errorObject.Message)
		assert.Equal(t, expectedErrorObject.Error(), errorObject.Error())
	})

	t.Run("Validation error", func(t *testing.T) {
		errorObject := eobject.ValidationError("`owner_type` object is empty")
		expectedErrorObject := entity.ErrorObject{
			Code:    "ERR1442",
			Message: fmt.Sprintf("Validation error because of: %s", "`owner_type` object is empty"),
		}

		assert.Equal(t, expectedErrorObject.Code, errorObject.Code)
		assert.Equal(t, expectedErrorObject.Message, errorObject.Message)
		assert.Equal(t, expectedErrorObject.Error(), errorObject.Error())
	})
}
