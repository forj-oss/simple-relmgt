package core

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// Config is the top Configuration object.
type Config struct {
	yaml yamlConfig
	file string
}

// NewConfig creates a Config object
func NewConfig(file string) (ret *Config) {
	ret = new(Config)

	ret.file = file
	return
}

// Load the configuration file
func (c *Config) Load() (err error) {
	if c == nil {
		return errors.New("Config object is nil. Unable to load")
	}

	var fd *os.File
	fd, err = os.Open(c.file)
	if err != nil {
		return fmt.Errorf("Unable to load '%s'. %s", c.file, err)
	}

	decoder := yaml.NewDecoder(fd)

	if decoder == nil {
		return fmt.Errorf("Unable to load '%s'. yaml decoder object not created", c.file)
	}

	err = decoder.Decode(&c.yaml)
	if err != nil {
		return fmt.Errorf("Unable to read yaml file '%s'. %s", c.file, err)
	}
	return
}
