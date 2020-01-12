package query

// oauth_applications
var PaginateOauthApplication = "select * from oauth_applications limit ?, ?"
var CountOauthApplication = "select count(*) as total from oauth_applications where revoked_at is null"
var SelectOneOauthApplication = "select * from oauth_applications where id = ?"
var SelectOneByUIDandSecret = "select * from oauth_applications where client_uid = ? and client_secret = ? limit 1"
var InsertOauthApplication = "insert into oauth_applications (owner_id, owner_type, description, scopes, client_uid, client_secret, revoked_at, created_at, updated_at) values(?, ?, ?, ?, ?, ?, null, now(), now())"

// oauth_access_tokens
var InsertOauthAccessToken = "insert into oauth_access_tokens (oauth_application_id, resource_owner_id, token, scopes, expires_in, created_at, revoked_at) values(?, ?, ?, ?, ?, now(), null)"
var SelectOneOauthAccessToken = "select * from oauth_access_tokens where id = ? limit 1"
