package main

type WordlistSpec struct {
	Sources []string `yaml:"sources"`
	Filters []string `yaml:"filters"`
}

type GenerationConfig struct {
	First  WordlistSpec `yaml:"first"`
	Second WordlistSpec `yaml:"second"`
}
