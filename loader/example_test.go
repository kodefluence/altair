package loader_test

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
