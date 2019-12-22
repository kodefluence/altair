package query_test

import (
	"testing"

	"github.com/codefluence-x/altair/query"
	"gotest.tools/assert"
)

func TestQuery(t *testing.T) {
	assert.Equal(t, "select * from oauth_applications limit ?, ?", query.PaginateOauthApplication)
	assert.Equal(t, "select count(*) from oauth_applications where revoked_at is null", query.CountOauthApplication)
	assert.Equal(t, "select * from oauth_applications where id = ?", query.SelectOneOauthApplication)
	assert.Equal(t, "insert into oauth_applications (owner_id, description, scopes, client_uid, client_secret, revoked_at, created_at, updated_at) values(?, ?, ?, ?, ?, null, now(), now())", query.InsertOauthApplication)
}
