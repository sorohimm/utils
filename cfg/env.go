package cfg

import (
	"fmt"
	"os"
	"strings"
)

// loadEnv reads the contents of the specified file and sets the environment variables accordingly.
// It takes a string argument `file` representing the path to the file to be loaded.
// It returns an error if there is any issue with reading the file or setting the environment variables.
// The file should be in a key=value format, with each line representing a separate key-value pair.
// Comments (lines starting with '#' or ';') are ignored.
// An error is returned if a line does not contain the expected key-value format.
// If an environment variable with the same key already exists, its value is not overridden.
// Example Usage:
// err := loadEnv(".env")
//
//	if err != nil {
//	   // handle error
//	}
func loadEnv(file string) error {
	fileContents, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("reading file %s error: %w", file, err)
	}

	lines := strings.Split(string(fileContents), "\n")
	for i, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		if isSkippable(trimmedLine) {
			continue
		}

		indexOfEqual := strings.Index(trimmedLine, "=")
		if indexOfEqual == -1 {
			return fmt.Errorf("file %s contains not expected line %s", file, trimmedLine)
		}

		keyValuePair := strings.SplitN(trimmedLine, "=", 2)
		if len(keyValuePair) != 2 {
			return fmt.Errorf("line at %v is %s but expected to be key=value", i, trimmedLine)
		}

		key, value := strings.TrimSpace(keyValuePair[0]), strings.TrimSpace(keyValuePair[1])
		currentVal := os.Getenv(key)
		if currentVal != "" {
			continue
		}

		err = os.Setenv(key, value)
		if err != nil {
			return fmt.Errorf("setting env %s=%s error: %w", key, value, err)
		}
	}

	return nil
}

// isSkippable check whether a line is skippable
func isSkippable(line string) bool {
	return line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";")
}
