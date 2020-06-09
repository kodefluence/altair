package entity

type Plugin struct {
	Plugin string `yaml:"plugin"`
	Raw    []byte `yaml:"-"`
}
