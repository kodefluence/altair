package cfg_test

var DatabaseConfigNormalScenario = `
oauth_database:
  driver: mysql
  database: {{ env "DATABASE_NAME_DB_CONFIG_NORMAL_SCENARIO" }}
  username: {{ env "DATABASE_USERNAME_DB_CONFIG_NORMAL_SCENARIO" }}
  password: {{ env "DATABASE_PASSWORD_DB_CONFIG_NORMAL_SCENARIO" }}
  host:     {{ env "DATABASE_HOST_DB_CONFIG_NORMAL_SCENARIO" }}
  port:     {{ env "DATABASE_PORT_DB_CONFIG_NORMAL_SCENARIO" }}
  connection_max_lifetime: 120s
  max_iddle_connection: 100
  max_open_connection: 100`

var DatabaseConfigWithNotFoundENV = `
  oauth_database:
    driver: mysql
    database: {{ env "DATABASE_NAME_DB_CONFIG_ENV_CONFIG_NOT_FOUND" }}
    username: {{ env "DATABASE_USERNAME_DB_CONFIG_ENV_CONFIG_NOT_FOUND" }}
    password: {{ env "DATABASE_PASSWORD_DB_CONFIG_ENV_CONFIG_NOT_FOUND" }}
    host:     {{ env "DATABASE_HOST_DB_CONFIG_ENV_CONFIG_NOT_FOUND" }}
    port:     {{ env "DATABASE_PORT_DB_CONFIG_ENV_CONFIG_NOT_FOUND" }}
    connection_max_lifetime: {{ env "NOT_FOUND_ENV" }}
    max_iddle_connection: 100
    max_open_connection: 100`

var DatabaseConfigNormalScenarioWithTwoValue = `
oauth_database:
  driver: mysql
  database: {{ env "DATABASE_NAME_TWO_SCENARIO_1" }}
  username: {{ env "DATABASE_USERNAME_TWO_SCENARIO_1" }}
  password: {{ env "DATABASE_PASSWORD_TWO_SCENARIO_1" }}
  host:     {{ env "DATABASE_HOST_TWO_SCENARIO_1" }}
  port:     {{ env "DATABASE_PORT_TWO_SCENARIO_1" }}
  connection_max_lifetime: 120s
  max_iddle_connection: 100
  max_open_connection: 100

other_database:
  driver: mysql
  database: {{ env "DATABASE_NAME_TWO_SCENARIO_2" }}
  username: {{ env "DATABASE_USERNAME_TWO_SCENARIO_2" }}
  password: {{ env "DATABASE_PASSWORD_TWO_SCENARIO_2" }}
  host:     {{ env "DATABASE_HOST_TWO_SCENARIO_2" }}
  port:     {{ env "DATABASE_PORT_TWO_SCENARIO_2" }}
  connection_max_lifetime: 120s
  max_iddle_connection: 100
  max_open_connection: 100`

var DatabaseConfigMissingDriver = `
oauth_database:
  database: some_database
  username: some_username
  password: some_password
  host:     localhost
  port:     3306
  connection_max_lifetime: 120s
  max_iddle_connection: 100
  max_open_connection: 100`

var DatabaseConfigInvalidDriver = `
oauth_database:
  driver: postgre
  database: some_database
  username: some_username
  password: some_password
  host:     localhost
  port:     3306
  connection_max_lifetime: 120s
  max_iddle_connection: 100
  max_open_connection: 100`

var DatabaseConfigMYSQLEmptyDatabase = `
oauth_database:
  driver: mysql
  database: ""
  username: some_username
  password: some_password
  host:     localhost
  port:     3306
  connection_max_lifetime: 120s
  max_iddle_connection: 100
  max_open_connection: 100`

var DatabaseConfigMYSQLEmptyUsername = `
oauth_database:
  driver: mysql
  database: some_database
  username: ""
  password: some_password
  host:     localhost
  port:     3306
  connection_max_lifetime: 120s
  max_iddle_connection: 100
  max_open_connection: 100`

var DatabaseConfigMYSQLEmptyHost = `
oauth_database:
  driver: mysql
  database: some_database
  username: some_username
  password: some_password
  host:     ""
  port:     3306
  connection_max_lifetime: 120s
  max_iddle_connection: 100
  max_open_connection: 100`

var DatabaseConfigMYSQLEmptyPort = `
oauth_database:
  driver: mysql
  database: some_database
  username: some_username
  password: some_password
  host:     localhost
  connection_max_lifetime: 120s
  max_iddle_connection: 100
  max_open_connection: 100`

var DatabaseConfigMYSQLEmptyConnectionMaxLifetime = `
oauth_database:
  driver: mysql
  database: some_database
  username: some_username
  password: some_password
  host:     localhost
  port:     3306
  max_iddle_connection: 100
  max_open_connection: 100`

var DatabaseConfigMYSQLEmptyMaxIddleConnection = `
oauth_database:
  driver: mysql
  database: some_database
  username: some_username
  password: some_password
  host:     localhost
  port:     3306
  connection_max_lifetime: 120s
  max_open_connection: 100`

var DatabaseConfigMYSQLEmptyMaxOpenConnection = `
oauth_database:
  driver: mysql
  database: some_database
  username: some_username
  password: some_password
  host:     localhost
  port:     3306
  connection_max_lifetime: 120s
  max_iddle_connection: 100`

var DatabaseConfigInvalidYaml = `
asdasd asda
a  asaa
-12 -2`

var DatabaseConfigV1Envelope = `
version: "1.0"
instances:
  oauth_database:
    driver: mysql
    database: {{ env "DATABASE_NAME_V1_ENVELOPE" }}
    username: {{ env "DATABASE_USERNAME_V1_ENVELOPE" }}
    password: {{ env "DATABASE_PASSWORD_V1_ENVELOPE" }}
    host:     {{ env "DATABASE_HOST_V1_ENVELOPE" }}
    port:     {{ env "DATABASE_PORT_V1_ENVELOPE" }}
    connection_max_lifetime: 120s
    max_idle_connection: 100
    max_open_connection: 100`

var DatabaseConfigV1EnvelopeMissingInstances = `
version: "1.0"`

var DatabaseConfigUnsupportedVersion = `
version: "9.9"
instances:
  oauth_database:
    driver: mysql
    database: x
    username: y
    host: z`

var DatabaseConfigIdleConnectionNewSpelling = `
version: "1.0"
instances:
  oauth_database:
    driver: mysql
    database: x
    username: y
    password: p
    host:     localhost
    port:     3306
    connection_max_lifetime: 120s
    max_idle_connection: 42
    max_open_connection: 100`

var DatabaseConfigInvalidTemplateFormatting = `
oauth_database:
  driver: mysql
  database: some_database
  username: some_username
  password: some_password
  host:     localhost
  port:     3306
  connection_max_lifetime: {}}}{}{{}}A{!@}
  max_iddle_connection: 100
  max_open_connection: 100`

var AppConfigNormal = `
version: 1.0
authorization:
  username: altair
  password: secret
plugins:
  - oauth
  - metric`

var AppConfigNormalWithoutVersion = `
version:
authorization:
  username: altair
  password: secret
plugins:
  - oauth
  - metric`

var AppConfigWithCustomProxyHost = `
version: 1.0
proxy_host: www.altair.id
authorization:
  username: altair
  password: secret
plugins:
  - oauth`

var AppConfigWithCustomPort = `
version: 1.0
port: 7001
authorization:
  username: altair
  password: secret
plugins:
  - oauth`

var AppConfigAuthUsernameEmpty = `
version: 1.0
authorization:
  password: secret
plugins:
  - oauth`

var AppConfigAuthPasswordEmpty = `
version: 1.0
authorization:
  username: altair
plugins:
  - oauth`

var AppConfigWithInvalidCustomPort = `
version: 1.0
port: asd
authorization:
  username: altair
  password: secret
plugins:
  - oauth`

var AppConfigUnmarshalError = `
ASd:
1231
AS
1
23
1
231
Aplugins:
- oauth
  - x`

var AppConfigTemplateError = `
plugins:
  - oauth
  - {{ } }} {}{} {}{`

var PluginConfigNormal1 = `
plugin: oauth
version: "1.0"
config:
  database: main_database
  access_token_timeout: 24h
  authorization_code_timeout: 24h`

var PluginConfigNormal2 = `
plugin: cache
version: "1.0"`

var PluginConfigMissingVersion = `
plugin: cache`

var PluginConfigMissingName = `
version: "1.0"`

var PluginConfigYamlUnmarshalError = `
asda

asd
1

23123;as;d_1231
`

var PluginConfigTemplateParsingError = `
plugin: oauth
version: "1.0"
config:
  database: {}}{}{{}{}}{}{}{{{{{}{}{}}}}}
  access_token_timeout: 24h
  authorization_code_timeout: 24h`
