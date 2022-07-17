package query_test

import (
	"testing"

	"github.com/kodefluence/altair/provider/plugin/oauth/query"
	"gotest.tools/assert"
)

func TestQuery(t *testing.T) {
	assert.Equal(t, "select id, owner_id, owner_type, description, scopes, client_uid, client_secret, revoked_at, created_at, updated_at from oauth_applications limit ?, ?", query.PaginateOauthApplication)
	assert.Equal(t, "select count(*) as total from oauth_applications where revoked_at is null", query.CountOauthApplication)
	assert.Equal(t, "select id, owner_id, owner_type, description, scopes, client_uid, client_secret, revoked_at, created_at, updated_at from oauth_applications where id = ?", query.SelectOneOauthApplication)
	assert.Equal(t, "select id, owner_id, owner_type, description, scopes, client_uid, client_secret, revoked_at, created_at, updated_at from oauth_applications where client_uid = ? and client_secret = ? limit 1", query.SelectOneByUIDandSecret)
	assert.Equal(t, "insert into oauth_applications (owner_id, owner_type, description, scopes, client_uid, client_secret, revoked_at, created_at, updated_at) values(?, ?, ?, ?, ?, ?, null, now(), now())", query.InsertOauthApplication)
	assert.Equal(t, "update oauth_applications set description = ?, scopes = ?, updated_at = now() where id = ?", query.UpdateOauthApplication)

	assert.Equal(t, "insert into oauth_access_tokens (oauth_application_id, resource_owner_id, token, scopes, expires_in, created_at, revoked_at) values(?, ?, ?, ?, ?, now(), null)", query.InsertOauthAccessToken)
	assert.Equal(t, "select id, oauth_application_id, resource_owner_id, token, scopes, expires_in, created_at, revoked_at from oauth_access_tokens where id = ? and revoked_at is null limit 1", query.SelectOneOauthAccessToken)
	assert.Equal(t, "select id, oauth_application_id, resource_owner_id, token, scopes, expires_in, created_at, revoked_at from oauth_access_tokens where token = ? and revoked_at is null limit 1", query.SelectOneOauthAccessTokenByToken)
	assert.Equal(t, "update oauth_access_tokens set revoked_at = now() where token = ?", query.RevokeAccessToken)

	assert.Equal(t, "insert into oauth_access_grants (oauth_application_id, resource_owner_id, scopes, code, redirect_uri, expires_in, created_at, revoked_at) values(?, ?, ?, ?, ?, ?, now(), null)", query.InsertOauthAccessGrant)
	assert.Equal(t, "select id, oauth_application_id, resource_owner_id, scopes, code, redirect_uri, expires_in, created_at, revoked_at from oauth_access_grants where id = ? limit 1", query.SelectOneOauthAccessGrant)
	assert.Equal(t, "select id, oauth_application_id, resource_owner_id, scopes, code, redirect_uri, expires_in, created_at, revoked_at from oauth_access_grants where code = ? limit 1", query.SelectOneOauthAccessGrantByCode)
	assert.Equal(t, "update oauth_access_grants set revoked_at = now() where code = ? limit 1", query.RevokeAuthorizationCode)

	assert.Equal(t, "insert into oauth_refresh_tokens (oauth_access_token_id, token, expires_in, created_at, revoked_at) values(?, ?, ?, now(), null)", query.InsertOauthRefreshToken)
	assert.Equal(t, "update oauth_refresh_tokens set revoked_at = now() where token = ?", query.RevokeRefreshToken)
	assert.Equal(t, "select id, oauth_access_token_id, token, expires_in, created_at, revoked_at from oauth_refresh_tokens where id = ? limit 1", query.SelectOneOauthRefreshToken)
	assert.Equal(t, "select id, oauth_access_token_id, token, expires_in, created_at, revoked_at from oauth_refresh_tokens where token = ? and revoked_at is null limit 1", query.SelectOneOauthRefreshTokenByToken)
}
