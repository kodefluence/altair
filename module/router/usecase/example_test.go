package usecase_test

var ExampleRoutesGracefully = `
name: users
auth: oauth
prefix: /users
host: {{ env "EXAMPLE_USERS_SERVICE_HOST" }}
path:
  /me: {}
  /:id: {}
`

var ExampleRoutesWithNoAuth = `
name: users
prefix: /users
host: {{ env "EXAMPLE_USERS_SERVICE_HOST" }}
path:
  /me: {}
  /:id: {}
`

var ExampleRoutesYamlError = `
name: 1
auth: 2
prefix: /users
host: {{ env "EXAMPLE_USERS_SERVICE_HOST" }}
this one make error
path:
  /me: {}
  /:id: {}
`

var ExampleTemplateParsingError = `
{{name: 1
auth: 2
prefix: /users
host: {.envasadasd}
path:
  /me: {}
  /:id: {}}}
`
var ExampleTemplateExecutionError = `
name: users
auth: execution_error
prefix: /users
host: {{ .env .x }}
path:
  /me: {}
  /:id: {}
`
