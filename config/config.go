package config

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/btm6084/godb"
	"github.com/btm6084/gojson"
	log "github.com/sirupsen/logrus"
)

var (
	rawConfig []byte

	// config is considered to be a singleton across an entire program.
	// Configuration that is expected to be overridden by environment variables should be at the root level
	// and be simple values. Secrets should (probably) come from the environment and not be comitted with code.
	config *gojson.JSONReader

	// env maps a BaseConfig key onto the Environmental Key that overrides it.
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
	env map[string]string

	// updateStopSignal stops the updater when written to. Only a single updater is allowed to run at a time.
	updateStopSignal = make(chan bool, 1)
	updating         bool

	defaultTimeout = 10 * time.Second

	ErrNoConfig     = errors.New("no config available. did you call init?")
	ErrNotValidJSON = errors.New("base config must be valid json")
)

// fetcher is used to retrieve override configs and a response to a query is expected to return JSON.
// If fetcher is nil, that functionality is disabled, and only baseConfig and envMap will be used.
//
//
// BaseConfig is the default configuration.
// Configuration that is expected to be overridden by environment variables should be at the root level
// and be simple values. Secrets should (probably) come from the environment and not be comitted with code.
//
// EnvMap maps an environment key onto the key in basic config to override.
// Environmental configurations must be strings in the BaseConfig, though you can use the Get*
// to convert to a desired type.
// Environmental configs are those things which are:
//     * Vital to the functionality of the service; or
//     * Secret in nature
// ex:
// "HOST":              "host"
// "PORT":              "port"
// "TIMEOUT":           "defaultTimeout"
// "ENV":               "env"
// "DB_CONNECT_STRING": "dbConnectString"
func Init(f godb.JSONFetcher, configFetchPath string, baseConfig []byte, envMap map[string]string, updateFrequency time.Duration) error {
	if !gojson.IsJSON(baseConfig) {
		return ErrNotValidJSON
	}

	rawConfig = baseConfig

	if updating {
		updateStopSignal <- true
		updating = false
	}

	err := update(f, configFetchPath, envMap)
	if err != nil {
		return err
	}

	if updateFrequency > 0 {
		go updater(f, configFetchPath, envMap, updateFrequency)
	}

	return nil
}

// SetConfig specifies a configuration value with no updating.
func SetConfig(b []byte) error {
	StopUpdates()

	rawConfig = b

	var err error
	config, err = gojson.NewJSONReader(rawConfig)
	if err != nil {
		return err
	}

	return nil
}

// StopUpdates sends a stop signal if there is a currently running updater.
func StopUpdates() {
	if updating {
		updateStopSignal <- true
		updating = false
	}
}

func updater(f godb.JSONFetcher, configFetchPath string, envMap map[string]string, updateFrequency time.Duration) error {
	updating = true
	ticker := time.NewTicker(updateFrequency)

	for {
		select {
		case <-updateStopSignal:
			return nil
		case <-ticker.C:
			err := update(f, configFetchPath, envMap)
			if err != nil {
				log.WithFields(log.Fields{
					"error":   err,
					"package": "github.com/btm6084/utilities/config",
				}).Error("configuration update failed")
			}
		}
	}
}

func update(f godb.JSONFetcher, configFetchPath string, envMap map[string]string) error {
	var fetched, envConfig []byte
	var err error

	fields := log.Fields{
		"package": "github.com/btm6084/utilities/config",
	}

	env = make(map[string]string)
	for envKey, cName := range envMap {
		envVal := os.Getenv(envKey)
		if envVal != "" {
			env[cName] = envVal
		}
	}

	// Merge base <- fetched <- env
	if f != nil {
		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
		defer cancel()

		fetched, err = f.FetchJSON(ctx, configFetchPath)
		if err != nil {
			fields["context"] = "update FetchJSON"
			log.WithFields(fields).Println(err)
		}
	}

	envConfig, err = json.Marshal(env)
	if err != nil {
		fields["context"] = "update marshal env"
		log.WithFields(fields).Println(err)
	}

	if !gojson.IsJSON(rawConfig) {
		rawConfig = []byte(`{}`)
	}

	if !gojson.IsJSON(fetched) {
		fetched = []byte(`{}`)
	}

	if !gojson.IsJSON(envConfig) {
		envConfig = []byte(`{}`)
	}

	c := merge(merge(rawConfig, fetched), envConfig)
	config, err = gojson.NewJSONReader(c)
	if err != nil {
		return err
	}

	rawConfig = c

	return nil
}

// merge b on top of a; get a back if it fails.
func merge(a, b []byte) []byte {
	out, err := gojson.MergeJSON(a, b)
	if err != nil {
		log.WithFields(log.Fields{"package": "github.com/btm6084/utilities/config", "context": "merge MergeJSON"}).Println(err)
		return a
	}

	return out
}

// GetString returns the value of the given configuration element at the specified key as a string,
// regardless of the original type at that key.
// Key is a valid JSON Path (https://goessner.net/articles/JsonPath/)
func GetString(key string) string {
	if config == nil || len(config.Keys) == 0 {
		log.Println(ErrNoConfig)
		return ""
	}

	return config.GetString(key)
}

// GetInt returns the value of the given configuration element at the specified key as an int,
// regardless of the original type at that key.
// Key is a valid JSON Path (https://goessner.net/articles/JsonPath/)
func GetInt(key string) int {
	if config == nil || len(config.Keys) == 0 {
		log.Println(ErrNoConfig)
		return 0
	}

	return config.GetInt(key)
}

// GetBool returns the value of the given configuration element at the specified key as a bool,
// regardless of the original type at that key.
// Key is a valid JSON Path (https://goessner.net/articles/JsonPath/)
func GetBool(key string) bool {
	if config == nil || len(config.Keys) == 0 {
		log.Println(ErrNoConfig)
		return false
	}

	return config.GetBool(key)
}

// GetFloat returns the value of the given configuration element at the specified key as a float64,
// regardless of the original type at that key.
// Key is a valid JSON Path (https://goessner.net/articles/JsonPath/)
func GetFloat(key string) float64 {
	if config == nil || len(config.Keys) == 0 {
		log.Println(ErrNoConfig)
		return 0.0
	}

	return config.GetFloat(key)
}

// GetMapStringBool returns the value of the given configuration element at the specified key as map[string]bool,
// regardless of the original type at that key.
// Key is a valid JSON Path (https://goessner.net/articles/JsonPath/)
func GetMapStringBool(key string) map[string]bool {
	if config == nil || len(config.Keys) == 0 {
		log.Println(ErrNoConfig)
		return nil
	}

	return config.Get(key).ToMapStringBool()
}

// GetMapStringString returns the value of the given configuration element at the specified key as map[string]string,
// regardless of the original type at that key.
// Key is a valid JSON Path (https://goessner.net/articles/JsonPath/)
func GetMapStringString(key string) map[string]string {
	if config == nil || len(config.Keys) == 0 {
		log.Println(ErrNoConfig)
		return nil
	}

	return config.Get(key).ToMapStringString()
}

// GetMapStringInt returns the value of the given configuration element at the specified key as map[string]int,
// regardless of the original type at that key.
// Key is a valid JSON Path (https://goessner.net/articles/JsonPath/)
func GetMapStringInt(key string) map[string]int {
	if config == nil || len(config.Keys) == 0 {
		log.Println(ErrNoConfig)
		return nil
	}

	return config.Get(key).ToMapStringInt()
}

// GetMapStringFloat returns the value of the given configuration element at the specified key as map[string]int,
// regardless of the original type at that key.
// Key is a valid JSON Path (https://goessner.net/articles/JsonPath/)
func GetMapStringFloat(key string) map[string]float64 {
	if config == nil || len(config.Keys) == 0 {
		log.Println(ErrNoConfig)
		return nil
	}

	return config.Get(key).ToMapStringFloat()
}

// GetMapStringInterface returns the value of the given configuration element at the specified key as map[string]interface{},
// regardless of the original type at that key.
// Key is a valid JSON Path (https://goessner.net/articles/JsonPath/)
func GetMapStringInterface(key string) map[string]interface{} {
	if config == nil || len(config.Keys) == 0 {
		log.Println(ErrNoConfig)
		return nil
	}

	return config.Get(key).ToMapStringInterface()
}

// Get returns the value of the given configuration element at the specified key as the JSON Type
// as it exists at that key.
// Key is a valid JSON Path (https://goessner.net/articles/JsonPath/)
func Get(key string) interface{} {
	if config == nil || len(config.Keys) == 0 {
		log.Println(ErrNoConfig)
		return nil
	}

	return config.GetInterface(key)
}
