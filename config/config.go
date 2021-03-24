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
	ErrNoConfig     = errors.New("no config available. did you call init?")
	ErrNotValidJSON = errors.New("base config must be valid json")
)

// merge b on top of a; get a back if it fails.
func merge(a, b []byte) []byte {
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

	updateStopSignal chan bool
	updating         bool
	updateFrequency  time.Duration
}

func NewLocalConfiguration(baseConfig []byte, envMap map[string]string) (*Configuration, error) {
	c := &Configuration{
		RawConfig:        baseConfig,
		EnvMap:           envMap,
		updating:         false,
		updateStopSignal: make(chan bool, 1),
		updateFrequency:  0,
	}

	err := c.update(nil, "")
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

	err := c.update(f, settingsPath)
	if err != nil {
		return nil, err
	}

	if c.updateFrequency > 0 {
		go c.updater(f, settingsPath)
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
				log.WithFields(log.Fields{
					"error":   err,
					"package": "github.com/btm6084/utilities/config",
				}).Error("configuration update failed")
			}
		}
	}
}

func (c *Configuration) update(f godb.JSONFetcher, path string) error {
	var fetched, envConfig []byte
	var err error

	fields := log.Fields{
		"package": "github.com/btm6084/utilities/config",
	}

	env := make(map[string]string)
	for envKey, cName := range c.EnvMap {
		envVal := os.Getenv(envKey)
		if envVal != "" {
			env[cName] = envVal
		}
	}

	// Merge base <- fetched <- env
	if f != nil {
		// It's assumed that the root url for the fetcher will provide the configuration data.
		fetched, err = f.FetchJSON(context.Background(), path)
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

	if !gojson.IsJSON(fetched) {
		fetched = []byte(`{}`)
	}

	if !gojson.IsJSON(envConfig) {
		envConfig = []byte(`{}`)
	}

	cfg := merge(merge(c.RawConfig, fetched), envConfig)
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
