package utils

import (
	"gopkg.in/yaml.v3"

	"os"
)

func ReadYamlFile(name string, cfg interface{}) error {
	data, err := os.ReadFile(name)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return err
	}
	return nil
}

// FileExists checks if a file exists and is not a directory
// before we try using it to prevent further errors.
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func CheckOrCreateFile(filename string) error {
	if FileExists(filename) {
		return nil
	}
	file, err := os.Create(filename)
	if err == nil {
		err = file.Close()
	}
	return err
}
