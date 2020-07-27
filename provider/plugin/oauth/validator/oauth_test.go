package validator_test

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"testing"

	"github.com/codefluence-x/altair/provider/plugin/oauth/entity"
	"github.com/codefluence-x/altair/provider/plugin/oauth/eobject"
	"github.com/codefluence-x/altair/provider/plugin/oauth/validator"
	"github.com/codefluence-x/altair/util"
	"github.com/stretchr/testify/assert"
)

func TestApplication(t *testing.T) {

	t.Run("ValidateApplication", func(t *testing.T) {
		t.Run("Given context and oauth application json data", func(t *testing.T) {
			t.Run("Return nil", func(t *testing.T) {
				data := entity.OauthApplicationJSON{
					OwnerID:     util.IntToPointer(1),
					OwnerType:   util.StringToPointer("confidential"),
					Description: util.StringToPointer("This is description"),
					Scopes:      util.StringToPointer("public users"),
				}
				applicationValidator := validator.Oauth()
				assert.Nil(t, applicationValidator.ValidateApplication(context.Background(), data))
			})
		})

		t.Run("Given context and oauth application json data with empty owner_type", func(t *testing.T) {
			t.Run("Return validation error", func(t *testing.T) {
				data := entity.OauthApplicationJSON{
					OwnerID:     util.IntToPointer(1),
					OwnerType:   nil,
					Description: util.StringToPointer("This is description"),
					Scopes:      util.StringToPointer("public users"),
				}

				expectedError := entity.Error{
					HttpStatus: http.StatusUnprocessableEntity,
					Errors:     eobject.Wrap(eobject.ValidationError("object `owner_type` is nil or not exists")),
				}

				applicationValidator := validator.Oauth()
				err := applicationValidator.ValidateApplication(context.Background(), data)

				assert.NotNil(t, err)
				assert.Equal(t, expectedError.HttpStatus, err.HttpStatus)
				assert.Equal(t, expectedError.Error(), err.Error())
				assert.Equal(t, expectedError.Errors, err.Errors)
			})
		})

		t.Run("Given context and oauth application json data with invalid owner_type", func(t *testing.T) {
			t.Run("Return validation error", func(t *testing.T) {
				data := entity.OauthApplicationJSON{
					OwnerID:     util.IntToPointer(1),
					OwnerType:   util.StringToPointer("external"),
					Description: util.StringToPointer("This is description"),
					Scopes:      util.StringToPointer("public users"),
				}

				expectedError := entity.Error{
					HttpStatus: http.StatusUnprocessableEntity,
					Errors:     eobject.Wrap(eobject.ValidationError("object `owner_type` must be either of `confidential` or `public`")),
				}

				applicationValidator := validator.Oauth()
				err := applicationValidator.ValidateApplication(context.Background(), data)

				assert.NotNil(t, err)
				assert.Equal(t, expectedError.HttpStatus, err.HttpStatus)
				assert.Equal(t, expectedError.Error(), err.Error())
				assert.Equal(t, expectedError.Errors, err.Errors)
			})
		})
	})

	t.Run("ValidateAuthorizationGrant", func(t *testing.T) {
		t.Run("Given context, authorization request and oauth application", func(t *testing.T) {
			t.Run("No scopes given", func(t *testing.T) {
				t.Run("Scopes is empty string", func(t *testing.T) {
					authorizationRequest := entity.AuthorizationRequestJSON{
						ResponseType:    util.StringToPointer("code"),
						ResourceOwnerID: util.IntToPointer(1),
						RedirectURI:     util.StringToPointer("www.github.com"),
						Scopes:          util.StringToPointer(""),
					}

					oauthApplication := entity.OauthApplication{
						Scopes: sql.NullString{
							String: "public users",
							Valid:  true,
						},
					}

					applicationValidator := validator.Oauth()
					err := applicationValidator.ValidateAuthorizationGrant(context.Background(), authorizationRequest, oauthApplication)
					assert.Nil(t, err)
				})

				t.Run("Scopes is nil", func(t *testing.T) {
					authorizationRequest := entity.AuthorizationRequestJSON{
						ResponseType:    util.StringToPointer("code"),
						ResourceOwnerID: util.IntToPointer(1),
						RedirectURI:     util.StringToPointer("www.github.com"),
						Scopes:          nil,
					}

					oauthApplication := entity.OauthApplication{
						Scopes: sql.NullString{
							String: "public users",
							Valid:  true,
						},
					}

					applicationValidator := validator.Oauth()
					err := applicationValidator.ValidateAuthorizationGrant(context.Background(), authorizationRequest, oauthApplication)
					assert.Nil(t, err)
				})
			})

			t.Run("Request scopes is unavailable in oauth application", func(t *testing.T) {
				authorizationRequest := entity.AuthorizationRequestJSON{
					ResponseType:    util.StringToPointer("code"),
					ResourceOwnerID: util.IntToPointer(1),
					RedirectURI:     util.StringToPointer("www.github.com"),
					Scopes:          util.StringToPointer("public users stores"),
				}

				oauthApplication := entity.OauthApplication{
					Scopes: sql.NullString{
						String: "public users",
						Valid:  true,
					},
				}

				applicationValidator := validator.Oauth()
				err := applicationValidator.ValidateAuthorizationGrant(context.Background(), authorizationRequest, oauthApplication)

				expectedErr := &entity.Error{
					HttpStatus: http.StatusForbidden,
					Errors:     eobject.Wrap(eobject.ForbiddenError(context.Background(), "application", fmt.Sprintf("your requested scopes `(%v)` is not exists in application", []string{"stores"}))),
				}

				assert.NotNil(t, err)
				assert.Equal(t, expectedErr.Error(), err.Error())
				assert.Equal(t, expectedErr.HttpStatus, err.HttpStatus)
				assert.Equal(t, expectedErr.Errors, err.Errors)
			})

			t.Run("Request scopes is available in oauth application", func(t *testing.T) {
				authorizationRequest := entity.AuthorizationRequestJSON{
					ResponseType:    util.StringToPointer("code"),
					ResourceOwnerID: util.IntToPointer(1),
					RedirectURI:     util.StringToPointer("www.github.com"),
					Scopes:          util.StringToPointer("public users"),
				}

				oauthApplication := entity.OauthApplication{
					Scopes: sql.NullString{
						String: "public users",
						Valid:  true,
					},
				}

				applicationValidator := validator.Oauth()
				err := applicationValidator.ValidateAuthorizationGrant(context.Background(), authorizationRequest, oauthApplication)
				assert.Nil(t, err)
			})
		})

		t.Run("Given context, authorization request without response type, resource owner id, redirect_uri and oauth application", func(t *testing.T) {
			t.Run("Return unprocessable entity", func(t *testing.T) {
				authorizationRequest := entity.AuthorizationRequestJSON{
					Scopes: util.StringToPointer(""),
				}

				oauthApplication := entity.OauthApplication{
					Scopes: sql.NullString{
						String: "public users",
						Valid:  true,
					},
				}

				expectedErr := &entity.Error{
					HttpStatus: http.StatusUnprocessableEntity,
					Errors: eobject.Wrap(
						eobject.ValidationError(`response_type can't be empty`),
						eobject.ValidationError(`resource_owner_id can't be empty`),
						eobject.ValidationError(`redirect_uri can't be empty`),
					),
				}

				applicationValidator := validator.Oauth()
				err := applicationValidator.ValidateAuthorizationGrant(context.Background(), authorizationRequest, oauthApplication)
				assert.NotNil(t, err)
				assert.Equal(t, expectedErr.Error(), err.Error())
				assert.Equal(t, expectedErr.HttpStatus, err.HttpStatus)
				assert.Equal(t, expectedErr.Errors, err.Errors)
			})
		})

		t.Run("Given context, authorization request with token and oauth application", func(t *testing.T) {
			t.Run("application is not confidential", func(t *testing.T) {
				authorizationRequest := entity.AuthorizationRequestJSON{
					ResponseType:    util.StringToPointer("token"),
					ResourceOwnerID: util.IntToPointer(1),
					RedirectURI:     util.StringToPointer("www.github.com"),
					Scopes:          util.StringToPointer(""),
				}

				oauthApplication := entity.OauthApplication{
					Scopes: sql.NullString{
						String: "public users",
						Valid:  true,
					},
				}

				ctx := context.Background()

				expectedErr := &entity.Error{
					HttpStatus: http.StatusForbidden,
					Errors: eobject.Wrap(
						eobject.ForbiddenError(ctx, "access_token", "your response type is not allowed in this application"),
					),
				}

				applicationValidator := validator.Oauth()
				err := applicationValidator.ValidateAuthorizationGrant(ctx, authorizationRequest, oauthApplication)
				assert.NotNil(t, err)
				assert.Equal(t, expectedErr.Error(), err.Error())
				assert.Equal(t, expectedErr.HttpStatus, err.HttpStatus)
				assert.Equal(t, expectedErr.Errors, err.Errors)
			})
		})
	})

	t.Run("ValidateTokenGrant", func(t *testing.T) {
		t.Run("Given context and access token request", func(t *testing.T) {
			ctx := context.Background()
			accessTokenRequest := entity.AccessTokenRequestJSON{
				ClientSecret: util.StringToPointer("client_secret"),
				ClientUID:    util.StringToPointer("client_uid"),
				Code:         util.StringToPointer("abcdef_123456"),
				GrantType:    util.StringToPointer("authorization_code"),
				RedirectURI:  util.StringToPointer("http://localhost:8000/oauth_redirect"),
			}

			t.Run("When access token request is valid", func(t *testing.T) {
				t.Run("Then return nil", func(t *testing.T) {
					applicationValidator := validator.Oauth()
					assert.Nil(t, applicationValidator.ValidateTokenGrant(ctx, accessTokenRequest))
				})
			})

			t.Run("When grant_type is nil", func(t *testing.T) {
				t.Run("Then return unprocessable entity", func(t *testing.T) {
					accessTokenRequest := entity.AccessTokenRequestJSON{
						ClientSecret: util.StringToPointer("client_secret"),
						ClientUID:    util.StringToPointer("client_uid"),
						Code:         util.StringToPointer("abcdef_123456"),
						// GrantType:    util.StringToPointer("authorization_code"),
						RedirectURI: util.StringToPointer("http://localhost:8000/oauth_redirect"),
					}

					expectedErr := &entity.Error{
						HttpStatus: http.StatusUnprocessableEntity,
						Errors: eobject.Wrap(
							eobject.ValidationError(`grant_type can't be empty`),
						),
					}

					applicationValidator := validator.Oauth()
					assert.Equal(t, expectedErr, applicationValidator.ValidateTokenGrant(ctx, accessTokenRequest))
				})
			})

			t.Run("When grant_type is not `authorization_code`", func(t *testing.T) {
				t.Run("Then return unprocessable entity", func(t *testing.T) {
					accessTokenRequest := entity.AccessTokenRequestJSON{
						ClientSecret: util.StringToPointer("client_secret"),
						ClientUID:    util.StringToPointer("client_uid"),
						Code:         util.StringToPointer("abcdef_123456"),
						GrantType:    util.StringToPointer("something_that_not_authorization_code"),
						RedirectURI:  util.StringToPointer("http://localhost:8000/oauth_redirect"),
					}

					expectedErr := &entity.Error{
						HttpStatus: http.StatusUnprocessableEntity,
						Errors: eobject.Wrap(
							eobject.ValidationError(`grant_type must be set to 'authorization_code'`),
						),
					}

					applicationValidator := validator.Oauth()
					assert.Equal(t, expectedErr, applicationValidator.ValidateTokenGrant(ctx, accessTokenRequest))
				})
			})

			t.Run("When code is nil", func(t *testing.T) {
				t.Run("Then return unprocessable entity", func(t *testing.T) {
					accessTokenRequest := entity.AccessTokenRequestJSON{
						ClientSecret: util.StringToPointer("client_secret"),
						ClientUID:    util.StringToPointer("client_uid"),
						// Code:         util.StringToPointer("abcdef_123456"),
						GrantType:   util.StringToPointer("authorization_code"),
						RedirectURI: util.StringToPointer("http://localhost:8000/oauth_redirect"),
					}

					expectedErr := &entity.Error{
						HttpStatus: http.StatusUnprocessableEntity,
						Errors: eobject.Wrap(
							eobject.ValidationError(`code can't be empty`),
						),
					}

					applicationValidator := validator.Oauth()
					assert.Equal(t, expectedErr, applicationValidator.ValidateTokenGrant(ctx, accessTokenRequest))
				})
			})

			t.Run("When redirect_uri is nil", func(t *testing.T) {
				t.Run("Then return unprocessable entity", func(t *testing.T) {
					accessTokenRequest := entity.AccessTokenRequestJSON{
						ClientSecret: util.StringToPointer("client_secret"),
						ClientUID:    util.StringToPointer("client_uid"),
						Code:         util.StringToPointer("abcdef_123456"),
						GrantType:    util.StringToPointer("authorization_code"),
						// RedirectURI: util.StringToPointer("http://localhost:8000/oauth_redirect"),
					}

					expectedErr := &entity.Error{
						HttpStatus: http.StatusUnprocessableEntity,
						Errors: eobject.Wrap(
							eobject.ValidationError(`redirect_uri can't be empty`),
						),
					}

					applicationValidator := validator.Oauth()
					assert.Equal(t, expectedErr, applicationValidator.ValidateTokenGrant(ctx, accessTokenRequest))
				})
			})
		})
	})
}
