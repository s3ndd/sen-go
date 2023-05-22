package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

// currentValues is the current config values
var currentValues atomic.Value

// Load loads and merges configuration from the runtime environment
func Load(locations ...string) {
	locations = append(locations, []string{".env.default", "etc/.env.default", ".env", "etc/.env"}...)

	location := os.Getenv("ENV_CONFIG_LOCATION")
	if location != "" {
		if !strings.HasPrefix(location, "/") {
			location = location + "/"
		}
		locations = append(locations, location)
	}

	values, err := NewEnvFromOsEnv()
	if err != nil {
		log.Fatal(err)
	}

	for i := len(locations) - 1; i >= 0; i-- {
		location = locations[i]
		log.Printf("Loading env from %s", location)
		env, err := NewEnvFromFile(location)
		if err != nil {
			if !os.IsNotExist(err) {
				log.Fatal(err)
			}
		}
		values = values.Merge(env)
	}

	currentValues.Store(values)
}

func ensureLoaded() {
	if currentValues.Load() == nil {
		Load()
	}
}

func readString(name string) string {
	return currentValues.Load().(Env).Get(name)
}

func Assert(name string) {
	ensureLoaded()
	if value := String(name, ""); value == "" {
		log.Fatal(fmt.Errorf("Config value %s is not set", name))
	}
}

// String returns a config value as a string.
// If the config value is not configured, it will return the default value.
func String(name string, defaultValue string) string {
	ensureLoaded()
	value := readString(name)
	if value == "" {
		return defaultValue
	}
	return value
}

// MustString will panic if the config value is not configured
func MustString(name string) string {
	Assert(name)
	return String(name, "")
}

// Strings returns a config value as a slice of strings.
func Strings(name string, defaultValue []string) []string {
	ensureLoaded()
	value := readString(name)
	if value == "" {
		return defaultValue
	}
	return strings.Split(value, ",")
}

// MustStrings will panic if the config value is not configured
func MustStrings(name string) []string {
	Assert(name)
	return Strings(name, []string{})
}

// Int returns a config value as an int.
// If the config value is not configured or cannot convert to an int, it will return the default value.
func Int(name string, defaultValue int) int {
	ensureLoaded()
	valueStr := readString(name)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

// MustInt will panic if the config value is not configured.
// If the config value cannot be converted to int, -1 will be returned.
func MustInt(name string) int {
	Assert(name)
	return Int(name, -1)
}

// Float returns a config value as a float64.
// If the config value is not configured or cannot convert to a float64, it will return the default value.
func Float(name string, defaultValue float64) float64 {
	ensureLoaded()
	valueStr := readString(name)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return defaultValue
	}

	return value
}

// MustFloat will panic if the config value is not configured.
// If the config value cannot be converted to float64, -1 will be returned.
func MustFloat(name string) float64 {
	Assert(name)
	return Float(name, -1)
}

// Duration returns a config value as a time.Duration.
// If the config value cannot be converted to time.Duration, the default value will be returned.
func Duration(name string, defaultValue time.Duration) time.Duration {
	ensureLoaded()
	valueStr := readString(name)
	if valueStr == "" {
		return defaultValue
	}

	duration, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}

	return duration
}

// MustDuration will panic if the config value is not configured.
// If the config value cannot be converted to time.Duration, 0 will be returned.
func MustDuration(name string) time.Duration {
	Assert(name)
	return Duration(name, 0)
}

// validBoolStrings is a map of valid bool strings
var validBoolStrings = map[string]bool{
	"false": false,
	"true":  true,
	"f":     false,
	"t":     true,
	"0":     false,
	"1":     true,
}

// Bool returns a config value as a bool.
// If the config value is not configured or cannot convert to a bool, it will return the default value.
func Bool(name string, defaultValue bool) bool {
	ensureLoaded()
	valueStr := readString(name)
	if valueStr == "" {
		return defaultValue
	}

	value, ok := validBoolStrings[strings.ToLower(valueStr)]
	if ok {
		return value
	}

	return defaultValue
}

// MustBool will panic if the config value is not configured.
// If the config value cannot be converted to bool, false will be returned.
func MustBool(name string) bool {
	Assert(name)
	return Bool(name, false)
}
