package eobject

import (
	"context"
	"fmt"

	"github.com/codefluence-x/altair/entity"
)

func Wrap(errorObject ...entity.ErrorObject) []entity.ErrorObject {
	return errorObject
}

func InternalServerError(ctx context.Context) entity.ErrorObject {
	return entity.ErrorObject{
		Code:    "ERR0500",
		Message: fmt.Sprintf("Something is not right, help us fix this problem. Contribute to https://github.com/codefluence-x/altair. Or help us by give this code '%v' to the admin of this site.", ctx.Value("track_id")),
	}
}

func BadRequestError(in string) entity.ErrorObject {
	return entity.ErrorObject{
		Code:    "ERR0400",
		Message: fmt.Sprintf("You've send malformed request in your `%s`", in),
	}
}
