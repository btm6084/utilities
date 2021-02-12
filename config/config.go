package config

import (
	"encoding/json"
	"os"

	"github.com/btm6084/gojson"
)

var (
	config *gojson.JSONReader

	// BaseConfig is the default configuration.
	// Configuration that is expected to be overridden by environment variables should be at the root level
	// and be simple values. Secrets should (probably) come from the environment and not be comitted with code.
	BaseConfig = map[string]interface{}{}

	// EnvMap maps a BaseConfig key onto the Environmental Key that overrides it.
	// Environmental configurations must be strings in the BaseConfig, though you can use the Get*
	// to convert to a desired type.
	// Environmental configs are those things which are:
	//     * Vital to the functionality of the service; or
	//     * Secret in nature
	// ex:
	// "host":              "HOST",
	// "port":              "PORT",
	// "defaultTimeout":    "TIMEOUT",
	// "env":               "ENV",
	// "dbConnectString":   "DB_CONNECT_STRING",
	EnvMap = map[string]string{}

	env map[string]string
)

func init() {
	c, err := json.Marshal(BaseConfig)
	if err != nil {
		panic(err)
	}

	env = make(map[string]string)
	for cName, envKey := range EnvMap {
		envVal := os.Getenv(envKey)
		if envVal != "" {
			env[cName] = envVal
		}
	}

	envJSON, err := json.Marshal(env)
	if err != nil {
		panic(err)
	}

	b, err := gojson.MergeJSON(c, envJSON)
	if err != nil {
		panic(err)
	}

	config, err = gojson.NewJSONReader(b)
	if err != nil {
		panic(err)
	}
}

// Merge takes a json string and merges it into the current config.
func Merge(in []byte) error {
	var original *gojson.JSONReader = config

	envJSON, err := json.Marshal(env)
	if err != nil {
		return err
	}

	c, err := json.Marshal(BaseConfig)
	if err != nil {
		return err
	}

	// Merge CurrentConfig <- IncomingConfig <- EnvironmentalConfig
	b, err := gojson.MergeJSON(c, in)
	if err != nil {
		return err
	}

	b, err = gojson.MergeJSON(b, envJSON)
	if err != nil {
		return err
	}

	config, err = gojson.NewJSONReader(b)
	if err != nil {
		*config = *original
		return err
	}

	return nil
}

// GetString returns the value of the given configuration element at the specified key as a string,
// regardless of the original type at that key.
// Key is a valid JSON Path (https://goessner.net/articles/JsonPath/)
func GetString(key string) string {
	return config.GetString(key)
}

// GetInt returns the value of the given configuration element at the specified key as an int,
// regardless of the original type at that key.
// Key is a valid JSON Path (https://goessner.net/articles/JsonPath/)
func GetInt(key string) int {
	return config.GetInt(key)
}

// GetBool returns the value of the given configuration element at the specified key as an bool,
// regardless of the original type at that key.
// Key is a valid JSON Path (https://goessner.net/articles/JsonPath/)
func GetBool(key string) bool {
	return config.GetBool(key)
}

// GetMapStringBool returns the value of the given configuration element at the specified key as map[string]bool,
// regardless of the original type at that key.
// Key is a valid JSON Path (https://goessner.net/articles/JsonPath/)
func GetMapStringBool(key string) map[string]bool {
	return config.Get(key).ToMapStringBool()
}

// GetMapStringString returns the value of the given configuration element at the specified key as map[string]string,
// regardless of the original type at that key.
// Key is a valid JSON Path (https://goessner.net/articles/JsonPath/)
func GetMapStringString(key string) map[string]string {
	return config.Get(key).ToMapStringString()
}

// GetMapStringInt returns the value of the given configuration element at the specified key as map[string]int,
// regardless of the original type at that key.
// Key is a valid JSON Path (https://goessner.net/articles/JsonPath/)
func GetMapStringInt(key string) map[string]int {
	return config.Get(key).ToMapStringInt()
}

// GetMapStringInterface returns the value of the given configuration element at the specified key as map[string]interface{},
// regardless of the original type at that key.
// Key is a valid JSON Path (https://goessner.net/articles/JsonPath/)
func GetMapStringInterface(key string) map[string]interface{} {
	return config.Get(key).ToMapStringInterface()
}

// Get returns the value of the given configuration element at the specified key as the JSON Type
// as it exists at that key.
// Key is a valid JSON Path (https://goessner.net/articles/JsonPath/)
func Get(key string) interface{} {
	config.GetInterface(key)

	return nil
}
