package envloader

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// EnvMap represents a set of environment variable key-value pairs.
type EnvMap map[string]string

// LoadFile reads a .env file and returns an EnvMap.
// Lines starting with '#' are treated as comments and skipped.
// Empty lines are also skipped.
func LoadFile(path string) (EnvMap, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening env file %q: %w", path, err)
	}
	defer f.Close()

	envMap := make(EnvMap)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid syntax at line %d: %q", lineNum, line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if key == "" {
			return nil, fmt.Errorf("empty key at line %d", lineNum)
		}

		envMap[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading env file %q: %w", path, err)
	}

	return envMap, nil
}
