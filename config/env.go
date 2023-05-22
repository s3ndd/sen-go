package config

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

// Env is a map of environment variables parsed from an env file
type Env map[string]string

// NewEnv loads envs from a file
func NewEnv() (Env, error) {
	return NewEnvFromFile(".env")
}

// NewEnvFromFile loads envs from a file
func NewEnvFromFile(name string) (Env, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	return load(file)
}

// NewEnvFromString loads envs from a string
func NewEnvFromOsEnv() (Env, error) {
	lines := os.Environ()
	if len(lines) == 0 {
		return Env{}, nil
	}
	var buf bytes.Buffer
	for _, line := range lines {
		if i := strings.IndexByte(line, '\n'); i >= 0 {
			line = line[:i]
		}
		_, err := buf.WriteString(line)
		if err != nil {
			return Env{}, err
		}
		buf.WriteByte('\n')
	}

	return load(&buf)
}

// Get return the value with the given key
func (e Env) Get(key string) string {
	return e[key]
}

// Merge returns a new env with the other env merged in
func (e Env) Merge(other Env) Env {
	ret := Env{}
	for key, value := range e {
		ret[key] = value
	}

	for key, value := range other {
		if _, present := e[key]; !present {
			ret[key] = value
		}
	}
	return ret
}

// Diff returns the keys that are different between the two envs
func (e Env) Diff(other Env) []string {
	diffs := map[string]bool{}
	for key, value := range e {
		if value != other.Get(key) {
			diffs[key] = true
		}
	}

	for key, value := range other {
		if seen := diffs[key]; !seen && value != e.Get(key) {
			diffs[key] = true
		}
	}

	diff := make([]string, 0, len(diffs))
	for name := range diffs {
		diff = append(diff, name)
	}

	return diff
}

// load loads the env from the reader
func load(reader io.Reader) (Env, error) {
	r := bufio.NewReader(reader)
	env := make(Env)
	for {
		line, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return env, err
		}
		lineStr := string(line)
		if len(lineStr) == 0 || strings.HasPrefix(lineStr, "#") {
			continue
		}
		key, val, err := parseEnvLine(lineStr)
		if err != nil {
			return env, err
		}

		env[key] = val
	}

	return env, nil
}

// parseEnvLine parses a line of env file
func parseEnvLine(line string) (string, string, error) {
	splits := strings.SplitN(line, "=", 2)
	if len(splits) < 2 {
		return "", "", fmt.Errorf("invalid env file due to missing delimiter'=' in line [%s]", line)
	}

	return strings.Trim(splits[0], " "), strings.Trim(splits[1], ` "'`), nil
}
