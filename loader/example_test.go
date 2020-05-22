package loader_test

var DatabaseConfigNormalScenario = `
oauth_database:
  driver: mysql
  database: {{ env "DATABASE_NAME_DB_CONFIG_NORMAL_SCENARIO" }}
  username: {{ env "DATABASE_USERNAME_DB_CONFIG_NORMAL_SCENARIO" }}
  password: {{ env "DATABASE_PASSWORD_DB_CONFIG_NORMAL_SCENARIO" }}
  migration_source: "file://migration"
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
    migration_source: "file://migration"
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
  migration_source: "file://migration"
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
  migration_source: "file://migration"
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
  migration_source: "file://migration"
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
  migration_source: "file://migration"
  host:     localhost
  port:     3306
  connection_max_lifetime: 120s
  max_iddle_connection: 100
  max_open_connection: 100`

var DatabaseConfigEmptyMigrationSource = `
oauth_database:
  driver: mysql
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
  migration_source: "file://migration"
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
  migration_source: "file://migration"
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
  migration_source: "file://migration"
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
  migration_source: "file://migration"
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
  migration_source: "file://migration"
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
  migration_source: "file://migration"
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
  migration_source: "file://migration"
  host:     localhost
  port:     3306
  connection_max_lifetime: 120s
  max_iddle_connection: 100`

var DatabaseConfigInvalidYaml = `
asdasd asda
a  asaa
-12 -2`

var DatabaseConfigInvalidTemplateFormatting = `
oauth_database:
  driver: mysql
  database: some_database
  username: some_username
  password: some_password
  migration_source: "file://migration"
  host:     localhost
  port:     3306
  connection_max_lifetime: {}}}{}{{}}A{!@}
  max_iddle_connection: 100
  max_open_connection: 100`

var AppConfigNormal = `
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
