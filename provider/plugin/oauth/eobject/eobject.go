package eobject

import (
	"context"
	"fmt"

	"github.com/kodefluence/altair/provider/plugin/oauth/entity"
)

func Wrap(errorObject ...entity.ErrorObject) []entity.ErrorObject {
	return errorObject
}

func InternalServerError(ctx context.Context) entity.ErrorObject {
	return entity.ErrorObject{
		Code:    "ERR0500",
		Message: fmt.Sprintf("Something is not right, help us fix this problem. Contribute to https://github.com/kodefluence/altair. Or help us by give this code '%v' to the admin of this site.", ctx.Value("request_id")),
	}
}

func BadRequestError(in string) entity.ErrorObject {
	return entity.ErrorObject{
		Code:    "ERR0400",
		Message: fmt.Sprintf("You've send malformed request in your `%s`", in),
	}
}

func NotFoundError(ctx context.Context, entityType string) entity.ErrorObject {
	return entity.ErrorObject{
		Code:    "ERR0404",
		Message: fmt.Sprintf("Resource of `%s` is not found, please report to admin of this site with this code `%v` if you think this is an error.", entityType, ctx.Value("request_id")),
	}
}

func UnauthorizedError() entity.ErrorObject {
	return entity.ErrorObject{
		Code:    "ERR0401",
		Message: fmt.Sprintf("You are unauthorized."),
	}
}

func ForbiddenError(ctx context.Context, entityType, reason string) entity.ErrorObject {
	return entity.ErrorObject{
		Code:    "ERR0403",
		Message: fmt.Sprintf("Resource of `%s` is forbidden to be accessed, because of: %s. Please report to admin of this site with this code `%v` if you think this is an error.", entityType, reason, ctx.Value("request_id")),
	}
}

func ValidationError(msg string) entity.ErrorObject {
	return entity.ErrorObject{
		Code:    "ERR1442",
		Message: fmt.Sprintf("Validation error because of: %s", msg),
	}
}
