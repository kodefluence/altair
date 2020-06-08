package core

type Migrator interface {
	Up() error
	Down() error
	Steps(steps int) error
	Close() (error, error)
}
