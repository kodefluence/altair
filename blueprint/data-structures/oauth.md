## `Oauth-Application` (object) - Oauth application data

+ id: 1 (number, required)
+ owner_id: 1 (number) - ID of the application owner, it can be user_id or something like that.
+ description: `Some cool application` (string) -  Application description.
+ scopes: `public user` (string) - Application scopes, separate with space.
+ client_uid: `fc19318dd13128ce14344d066510a982269c241b` (string, required) - Client unique id.
+ client_secret: `AURPHM_23c6834b1d353eabf976e524ed489c812ff86a7d.23c6834b1d353eabf976e524ed489c812ff86a7d` (string, required) - Client secret token.
+ revoked_at: `` (string) - Date of application revoked.
+ created_at: `2020-01-01T00:00:00Z` (string, required) - Data created at time.
+ updated_at: `2020-01-01T00:00:00Z` (string, required) - Data last updated at time.


## `Oauth-Application-Other` (object) - Oauth application data

+ id: 2 (number, required)
+ owner_id: 2 (number) - ID of the application owner, it can be user_id or something like that.
+ description: `Some cool other application` (string) -  Application description.
+ scopes: `public user` (string) - Application scopes, separate with space.
+ client_uid: `fc19318dd13128ce14344d066510a982269c241b` (string, required) - Client unique id.
+ client_secret: `AURPHM_23c6834b1d353eabf976e524ed489c812ff86a7d.23c6834b1d353eabf976e524ed489c812ff86a7d` (string, required) - Client secret token.
+ revoked_at: `` (string) - Date of application revoked.
+ created_at: `2020-01-01T00:00:00Z` (string, required) - Data created at time.
+ updated_at: `2020-01-01T00:00:00Z` (string, required) - Data last updated at time.

## `Oauth-Application-Create-Request` (object) - Oauth application create request data

+ owner_id: 1 (number) - ID of the application owner, it can be user_id or something like that.
+ description: `Some cool application` (string) -  Application description.
+ scopes: `public user` (string) - Application scopes, separate with space.