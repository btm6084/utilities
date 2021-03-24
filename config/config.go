package config

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/btm6084/godb"
	"github.com/btm6084/gojson"
	"github.com/btm6084/utilities/stack"
	log "github.com/sirupsen/logrus"
)

var (
	ErrNoConfig     = errors.New("no config available. did you call init?")
	ErrNotValidJSON = errors.New("base config must be valid json")
)

// merge b on top of a; get a back if it fails.
func merge(a, b []byte) []byte {
	aEmpty := gojson.IsEmptyObject(a)
	bEmpty := gojson.IsEmptyObject(b)

	if !aEmpty && bEmpty {
		return a
	}

	if aEmpty && !bEmpty {
		return b
	}

	out, err := gojson.MergeJSON(a, b)
	if err != nil {
		log.WithFields(log.Fields{"package": "github.com/btm6084/utilities/config", "context": "merge MergeJSON"}).Println(err)
		return a
	}
	return out
}

type Configuration struct {
	Config    *gojson.JSONReader
	RawConfig []byte
	EnvMap    map[string]string

	envConfig     []byte
	fetchedConfig []byte

	updateStopSignal chan bool
	updating         bool
	updateFrequency  time.Duration
}

func (c *Configuration) parseEnvConfig() ([]byte, error) {
	env := make(map[string]string)
	for envKey, cName := range c.EnvMap {
		envVal := os.Getenv(envKey)
		if envVal != "" {
			env[cName] = envVal
		}
	}

	envConfig, err := json.Marshal(env)
	if err != nil {
		return []byte(`{}`), err
	}

	return envConfig, nil
}

func NewLocalConfiguration(baseConfig []byte, envMap map[string]string) (*Configuration, error) {
	if !gojson.IsJSON(baseConfig) {
		return nil, ErrNotValidJSON
	}

	c := &Configuration{
		RawConfig:        baseConfig,
		EnvMap:           envMap,
		updating:         false,
		updateStopSignal: make(chan bool, 1),
		updateFrequency:  0,
	}

	var err error
	c.envConfig, err = c.parseEnvConfig()
	if err != nil {
		log.WithFields(stack.TraceFields()).Error(err)
	}

	c.RawConfig = merge(c.RawConfig, c.envConfig)
	reader, err := gojson.NewJSONReader(c.RawConfig)
	if err != nil {
		return nil, err
	}

	c.Config = reader
	return c, err
}

func NewRemoteConfiguration(baseConfig []byte, envMap map[string]string, f godb.JSONFetcher, settingsPath string, updateFrequency time.Duration) (*Configuration, error) {
	if !gojson.IsJSON(baseConfig) {
		return nil, ErrNotValidJSON
	}

	c := &Configuration{
		RawConfig:        baseConfig,
		EnvMap:           envMap,
		updating:         false,
		updateStopSignal: make(chan bool, 1),
		updateFrequency:  updateFrequency,
	}

	var err error
	c.envConfig, err = c.parseEnvConfig()
	if err != nil {
		log.WithFields(stack.TraceFields()).Error(err)
	}

	c.RawConfig = merge(c.RawConfig, c.envConfig)

	err = c.update(f, settingsPath)
	if err != nil && err != godb.ErrNotFound {
		return nil, err
	}

	if f != nil && c.updateFrequency > 0 {
		go c.updater(f, settingsPath)
	}

	if c.Config == nil {
		reader, err := gojson.NewJSONReader(c.RawConfig)
		if err != nil {
			return nil, err
		}

		c.Config = reader
	}

	return c, nil
}

func (c *Configuration) updater(f godb.JSONFetcher, path string) error {
	c.updating = true
	ticker := time.NewTicker(c.updateFrequency)

	for {
		select {
		case <-c.updateStopSignal:
			return nil
		case <-ticker.C:
			err := c.update(f, path)
			if err != nil {
				fields := stack.TraceFields()
				fields["error"] = err
				log.WithFields(fields).Error("configuration update failed")
			}
		}
	}
}

func (c *Configuration) update(f godb.JSONFetcher, path string) error {
	if f == nil {
		return nil
	}

	fetchedConfig, err := f.FetchJSON(context.Background(), path)
	if err != nil {
		return err
	}

	// Nothing to do here.
	if string(c.fetchedConfig) == string(fetchedConfig) {
		return nil
	}

	c.fetchedConfig = fetchedConfig

	// env must overwrite anything pulled from the remote settings
	cfg := merge(merge(c.RawConfig, c.fetchedConfig), c.envConfig)
	reader, err := gojson.NewJSONReader(cfg)
	if err != nil {
		return err
	}

	c.RawConfig = cfg
	c.Config = reader

	return nil
}

// StopUpdates sends a stop signal if there is a currently running updater.
func (c *Configuration) StopUpdates() {
	if c.updating {
		c.updateStopSignal <- true
		c.updating = false
	}
}
