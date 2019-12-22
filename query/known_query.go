package query

var PaginateOauthApplication = "select * from oauth_applications limit ?, ?"
var CountOauthApplication = "select count(*) from oauth_applicaions where revoked_at is null"
var SelectOneOauthApplication = "select * from oauth_applications where id = ?"
var InsertOauthApplication = "insert into oauth_applications (owner_id, description, scopes, client_uid, client_secret, revoked_at, created_at, updated_at) values(?, ?, ?, ?, ?, null, now(), now())"
