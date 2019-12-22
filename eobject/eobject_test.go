package eobject_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/eobject"
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
			Code:    "ERR500",
			Message: fmt.Sprintf("Something is not right, help us fix this problem. Contribute to https://github.com/codefluence-x/altair. Or help us by give this code '%v' to the admin of this site.", ctx.Value("track_id")),
		}

		assert.Equal(t, expectedErrorObject.Code, errorObject.Code)
		assert.Equal(t, expectedErrorObject.Message, errorObject.Message)
		assert.Equal(t, expectedErrorObject.Error(), errorObject.Error())
	})
}
