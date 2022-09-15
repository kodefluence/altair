package entity

type Plugin struct {
	Plugin  string `yaml:"plugin"`
	Version string `yaml:"version"`
	Raw     []byte `yaml:"-"`
}
