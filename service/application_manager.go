package service

import (
	"context"
	"net/http"

	"github.com/codefluence-x/altair/core"
	"github.com/codefluence-x/altair/entity"
	"github.com/codefluence-x/altair/eobject"
)

type applicationManager struct {
	formatter  core.OauthApplicationFormater
	oauthModel core.OauthApplicationModel
}

func ApplicationManager(formatter core.OauthApplicationFormater, oauthModel core.OauthApplicationModel) core.ApplicationManager {
	return &applicationManager{
		formatter:  formatter,
		oauthModel: oauthModel,
	}
}

func (am *applicationManager) List(ctx context.Context, offset, limit int) ([]entity.OauthApplicationJSON, int, *entity.Error) {
	oauthApplications, err := am.oauthModel.Paginate(ctx, offset, limit)
	if err != nil {
		return []entity.OauthApplicationJSON(nil), 0, &entity.Error{
			HttpStatus: http.StatusInternalServerError,
			Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
		}
	}

	total, err := am.oauthModel.Count(ctx)
	if err != nil {
		return []entity.OauthApplicationJSON(nil), 0, &entity.Error{
			HttpStatus: http.StatusInternalServerError,
			Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
		}
	}

	formattedResult := am.formatter.ApplicationList(ctx, oauthApplications)
	return formattedResult, total, nil
}
