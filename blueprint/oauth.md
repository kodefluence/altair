## Oauth Application [/oauth/applications]

These are the list of API of Oauth application in Altair.

### Get list oauth application [GET /oauth/applications{?offset,limit}]

Get a list of oauth application.

> BASIC AUTH

+ Parameters
  + offset (number, optional) - Pagination offset.
      + Default: `0`
  + limit (number, optional) - Pagination limit.
      + Default: `10`

+ Request (application/json)
  + Headers

            Authorization: Basic YWx0YWlyOmVhZ2xldGhhdGZseWludGhlYmx1ZXNreQ==

+ Response 200 (application/json)
  + Attributes (object)
      + data (array, required)
          + (object)
              + Include Oauth-Application
          + (object)
              + Include Oauth-Application-Other
      + meta (object, required)
          + offset: 0 (number, required)
          + limit: 10 (number, required)
          + total: 2 (number, required)

### Get oauth application by id [GET /oauth/applications/{id}]

Get an oauth application by id.

> BASIC AUTH

+ Parameters
  + id (number, required) - ID of oauth application.

+ Request (application/json)
  + Headers

            Authorization: Basic YWx0YWlyOmVhZ2xldGhhdGZseWludGhlYmx1ZXNreQ==

+ Response 200 (application/json)
  + Attributes (object)
      + data (object, required)
            + Include Oauth-Application

### Create oauth application [POST /oauth/applications]

Create new oauth application.

> BASIC AUTH

+ Request (application/json)
  + Headers

            Authorization: Basic YWx0YWlyOmVhZ2xldGhhdGZseWludGhlYmx1ZXNreQ==

  + Attributes (object)
      + Include Oauth-Application-Create-Request

+ Response 200 (application/json)
  + Attributes (object)
      + data (object, required)
          + Include Oauth-Application
