package usecase

// Create oauth application
// func (am *ApplicationManager) Create(ktx kontext.Context, e entity.OauthApplicationJSON) (entity.OauthApplicationJSON, jsonapi.Errors) {
// 	// 	if err := am.applicationValidator.ValidateApplication(ctx, e); err != nil {
// 	// 		log.Error().
// 	// 			Err(err).
// 	// 			Stack().
// 	// 			Interface("data", e).
// 	// 			Array("tags", zerolog.Arr().Str("service").Str("application_manager").Str("create").Str("validate_application")).
// 	// 			Msg("Got validation error from oauth application validator")
// 	// 		return entity.OauthApplicationJSON{}, err
// 	// 	}

// 	// 	id, err := am.oauthApplicationModel.Create(kontext.Fabricate(kontext.WithDefaultContext(ctx)), am.modelFormatter.OauthApplication(e), am.sqldb)
// 	// 	if err != nil {
// 	// 		log.Error().
// 	// 			Err(err).
// 	// 			Stack().
// 	// 			Interface("data", e).
// 	// 			Array("tags", zerolog.Arr().Str("service").Str("application_manager").Str("create").Str("model_create")).
// 	// 			Msg("Error when creating oauth application data")

// 	// 		return entity.OauthApplicationJSON{}, &entity.Error{
// 	// 			HttpStatus: http.StatusInternalServerError,
// 	// 			Errors:     eobject.Wrap(eobject.InternalServerError(ctx)),
// 	// 		}
// 	// 	}

// 	// 	return am.One(ctx, id)
// }
