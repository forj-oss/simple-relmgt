package core

import (
	"errors"
	"fmt"
	"os"

	"github.com/forj-oss/forjj-modules/trace"
	"gopkg.in/yaml.v2"
)

// Config is the top Configuration object.
type Config struct {
	Yaml YamlConfig
	file string
}

// NewConfig creates a Config object
func NewConfig(file string) (ret *Config) {
	ret = new(Config)

	ret.file = file
	return
}

// Filename return the configuration file name.
func (c Config) Filename() string {
	return c.file
}

// Load the configuration file
func (c *Config) Load() (loaded bool, err error) {
	if c == nil {
		err = errors.New("Config object is nil. Unable to load")
		return
	}

	var fd *os.File
	fd, err = os.Open(c.file)
	if err != nil {
		gotrace.Warning("%s not loaded. %s", c.file, err)
		return
	}

	decoder := yaml.NewDecoder(fd)

	if decoder == nil {
		err = fmt.Errorf("Unable to load '%s'. yaml decoder object not created", c.file)
		return
	}

	err = decoder.Decode(&c.Yaml)
	if err != nil {
		err = fmt.Errorf("Unable to read yaml file '%s'. %s", c.file, err)
		return
	}

	gotrace.Info("%s loaded.", c.file)
	loaded = true
	return
}
