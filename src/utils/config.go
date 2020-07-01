package utils

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Source struct {
	Name      string `yaml:"name"`
	RateLimit int8   `yaml:"ratelimit"`
}

type ConfigFile struct {
	Sources []Source `yaml:"sources"`
}

func GetConfig() ConfigFile {
	config := ConfigFile{}

	path, _ := filepath.Abs("../config.yaml")

	file, err := os.Open(path)
	FailOnError(err)

	defer file.Close()

	b, err := ioutil.ReadAll(file)

	data := string(b)

	err = yaml.Unmarshal([]byte(data), &config)
	FailOnError(err)

	return config
}
