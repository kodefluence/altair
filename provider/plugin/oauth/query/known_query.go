package query

// oauth_applications
const PaginateOauthApplication = "select * from oauth_applications limit ?, ?"
const CountOauthApplication = "select count(*) as total from oauth_applications where revoked_at is null"
const SelectOneOauthApplication = "select * from oauth_applications where id = ?"
const SelectOneByUIDandSecret = "select * from oauth_applications where client_uid = ? and client_secret = ? limit 1"
const InsertOauthApplication = "insert into oauth_applications (owner_id, owner_type, description, scopes, client_uid, client_secret, revoked_at, created_at, updated_at) values(?, ?, ?, ?, ?, ?, null, now(), now())"

// oauth_access_tokens
const InsertOauthAccessToken = "insert into oauth_access_tokens (oauth_application_id, resource_owner_id, token, scopes, expires_in, created_at, revoked_at) values(?, ?, ?, ?, ?, now(), null)"
const SelectOneOauthAccessToken = "select * from oauth_access_tokens where id = ? limit 1"
const SelectOneOauthAccessTokenByToken = "select id, oauth_application_id, resource_owner_id, token, scopes, expires_in, created_at, revoked_at from oauth_access_tokens where token = ? and revoked_at is null limit 1"
const RevokeAccessToken = "update oauth_access_tokens set revoked_at = now() where token = ?"

// oauth_access_grants
const InsertOauthAccessGrant = "insert into oauth_access_grants (oauth_application_id, resource_owner_id, scopes, code, redirect_uri, expires_in, created_at, revoked_at) values(?, ?, ?, ?, ?, ?, now(), null)"
const SelectOneOauthAccessGrant = "select * from oauth_access_grants where id = ? limit 1"
