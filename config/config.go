package config

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/btm6084/godb"
	"github.com/btm6084/gojson"
	log "github.com/sirupsen/logrus"
)

var (
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

	// fetcher is used to retrieve override configs and a response to a query is expected to return JSON.
	fetcher godb.JSONFetcher

	// updateStopSignal stops the updater when written to. Only a single updater is allowed to run at a time.
	updateStopSignal = make(chan bool, 1)
	updating         bool

	defaultTimeout = 10 * time.Second
)

// fetcher is used to retrieve override configs and a response to a query is expected to return JSON.
// If fetcher is nil, that functionality is disabled, and only baseConfig and envMap will be used.
//
//
// BaseConfig is the default configuration.
// Configuration that is expected to be overridden by environment variables should be at the root level
// and be simple values. Secrets should (probably) come from the environment and not be comitted with code.
//
// EnvMap maps a BaseConfig key onto the Environmental Key that overrides it.
// Environmental configurations must be strings in the BaseConfig, though you can use the Get*
// to convert to a desired type.
// Environmental configs are those things which are:
//     * Vital to the functionality of the service; or
//     * Secret in nature
// ex:
// "host":              "HOST"
// "port":              "PORT"
// "defaultTimeout":    "TIMEOUT"
// "env":               "ENV"
// "dbConnectString":   "DB_CONNECT_STRING"
func Init(f godb.JSONFetcher, configFetchPath string, baseConfig map[string]string, envMap map[string]string, updateFrequency time.Duration) error {
	if updating {
		updateStopSignal <- true
		updating = false
	}

	err := update(f, configFetchPath, baseConfig, envMap)
	if err != nil {
		return err
	}

	go updater(f, configFetchPath, baseConfig, envMap, updateFrequency)

	return nil
}

// StopUpdates sends a stop signal if there is a currently running updater.
func StopUpdates() {
	if updating {
		updateStopSignal <- true
		updating = false
	}
}

func updater(f godb.JSONFetcher, configFetchPath string, baseConfig map[string]string, envMap map[string]string, updateFrequency time.Duration) error {
	updating = true
	ticker := time.NewTicker(500 * time.Millisecond)

	for {
		select {
		case <-updateStopSignal:
			return nil
		case <-ticker.C:
			err := update(f, configFetchPath, baseConfig, envMap)
			if err != nil {
				log.WithFields(log.Fields{"error": err}).Error("configuration update failed")
			}
		}
	}
}

func update(f godb.JSONFetcher, configFetchPath string, baseConfig map[string]string, envMap map[string]string) error {
	var base, fetched, envConfig []byte
	var err error

	base, err = json.Marshal(baseConfig)
	if err != nil {
		return err
	}

	env = make(map[string]string)
	for cName, envKey := range envMap {
		envVal := os.Getenv(envKey)
		if envVal != "" {
			env[cName] = envVal
		}
	}

	// Merge base <- fetched <- env
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	if fetcher != nil {
		fetched, err = fetcher.FetchJSON(ctx, configFetchPath)
		if err != nil {
			log.Println(err)
		}
	}

	envConfig, err = json.Marshal(env)
	if err != nil {
		log.Println(err)
	}

	c := merge(merge(base, fetched), envConfig)
	config, err = gojson.NewJSONReader(c)
	if err != nil {
		return err
	}

	return nil
}

// merge b on top of a; get a back if it fails.
func merge(a, b []byte) []byte {
	out, err := gojson.MergeJSON(a, b)
	if err != nil {
		log.Println(err)
		return a
	}

	return out
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
