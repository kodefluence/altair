package entity

type MetricPlugin struct {
	Config struct {
		Provider string `yaml:"provider"`
	} `yaml:"config"`
}
