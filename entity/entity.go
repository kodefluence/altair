package entity

type RouteObject struct {
	Name   string                `yaml:"name"`
	Auth   string                `yaml:"auth"`
	Prefix string                `yaml:"prefix"`
	Host   string                `yaml:"host"`
	Path   map[string]RouterPath `yaml:"path"`
}

type RouterPath struct {
	Auth  string `yaml:"auth"`
	Scope string `yaml:"scope"`
}

func (r *RouterPath) GetAuth() string {
	return r.Auth
}

func (r *RouterPath) GetScope() string {
	return r.Scope
}
