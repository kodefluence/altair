package query

// PaginateOauthApplication oauth_applications table query
const PaginateOauthApplication = "select * from oauth_applications limit ?, ?"

// CountOauthApplication oauth_applications table query
const CountOauthApplication = "select count(*) as total from oauth_applications where revoked_at is null"

// SelectOneOauthApplication oauth_applications table query
const SelectOneOauthApplication = "select * from oauth_applications where id = ?"

// SelectOneByUIDandSecret oauth_applications table query
const SelectOneByUIDandSecret = "select * from oauth_applications where client_uid = ? and client_secret = ? limit 1"

// InsertOauthApplication oauth_applications table query
const InsertOauthApplication = "insert into oauth_applications (owner_id, owner_type, description, scopes, client_uid, client_secret, revoked_at, created_at, updated_at) values(?, ?, ?, ?, ?, ?, null, now(), now())"

// InsertOauthAccessToken oauth_access_tokens table query
const InsertOauthAccessToken = "insert into oauth_access_tokens (oauth_application_id, resource_owner_id, token, scopes, expires_in, created_at, revoked_at) values(?, ?, ?, ?, ?, now(), null)"

// SelectOneOauthAccessToken oauth_access_tokens table query
const SelectOneOauthAccessToken = "select * from oauth_access_tokens where id = ? limit 1"

// SelectOneOauthAccessTokenByToken oauth_access_tokens table query
const SelectOneOauthAccessTokenByToken = "select id, oauth_application_id, resource_owner_id, token, scopes, expires_in, created_at, revoked_at from oauth_access_tokens where token = ? and revoked_at is null limit 1"

// RevokeAccessToken oauth_access_tokens table query
const RevokeAccessToken = "update oauth_access_tokens set revoked_at = now() where token = ?"

// InsertOauthAccessGrant oauth_access_grants table query
const InsertOauthAccessGrant = "insert into oauth_access_grants (oauth_application_id, resource_owner_id, scopes, code, redirect_uri, expires_in, created_at, revoked_at) values(?, ?, ?, ?, ?, ?, now(), null)"

// SelectOneOauthAccessGrant oauth_access_grants table query
const SelectOneOauthAccessGrant = "select * from oauth_access_grants where id = ? limit 1"

// SelectOneOauthAccessGrantByCode oauth_access_grants table query
const SelectOneOauthAccessGrantByCode = "select * from oauth_access_grants where code = ?"
