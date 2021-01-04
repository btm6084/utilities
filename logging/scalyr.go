package logging

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	client = &http.Client{
		Timeout: 2 * time.Second,
	}
)

// ScalyrWriter will buffer all log writes for logging to scalyr, then
// tee all calls to the supplied io.Writer
type ScalyrWriter struct {
	// tee provides an io.Writer that receives all log entries, along with the entries written to scalyr.
	// popular options include os.Stdout and ioutil.Discard (/dev/null)
	tee        io.Writer
	buffer     bytes.Buffer
	url        string
	lock       sync.Mutex
	errorCount int
}

func (w *ScalyrWriter) Write(p []byte) (int, error) {
	// If we can't upload to Scalyr and flush the buffer, we stop accumulating
	// so that we don't overrun memory by infinitely buffering.
	if w.errorCount < 100 {
		w.lock.Lock()
		w.buffer.Write(p)
		w.lock.Unlock()
	}
	return w.tee.Write(p)
}

// CreateScalyrWriter will create an io.Writer which will tee all log writes
// between the supplied io.Writer and a buffer. You can then call
// go ScalyrWriter.Log with a supplied interval. On that interval, all collected
// logs will be sent to Scalyr
//
// URL Takes the following form, which requires the user to fill out the Host, Logfile, Access Token and Parser
// ex: https://www.scalyr.com/api/uploadLogs?host=my-host-name&logfile=myErrorLog&token=myScalyrToken&parser=goErrorLog
//
// Logs are POST'd in batches to Scalyr based on the interval provided to the Update function.
// Use logging.NoopWriter if you wish to ONLY log to Scalyr
func CreateScalyrWriter(tee io.Writer, url string) *ScalyrWriter {
	if tee == nil {
		tee = os.Stdout
	}

	return &ScalyrWriter{tee: tee, url: url}
}

// UpdateNow immediately uploads the collected logs to Scalyr.
// Interval should be provided in Milliseconds.
//
// Example usage:
// w := CreateScalyrWriter(os.Stdout, "https://www.scalyr.com/api/uploadLogs?host=ExampleService&logfile=AccessLog&token=ExampleToken")
// defer w.UpdateNow() // Update after leaving the current function.
func (w *ScalyrWriter) UpdateNow() {
	var buf []byte

	// Make a temporary copy of the buffer so we can release it asap.
	w.lock.Lock()

	if w.buffer.Len() > 0 {
		buf = make([]byte, w.buffer.Len())
		copy(buf, w.buffer.Bytes())
	}

	w.lock.Unlock()

	if len(buf) > 0 {
		start := time.Now()
		resp, err := client.Post(w.url, "text/plain", bytes.NewBuffer(buf))
		if err != nil {
			w.Println(err.Error())
			w.errorCount++
			return
		}
		w.Println("Scalyr Post Time:", time.Since(start).String(), "Scalyr Post Url:", w.url)

		r, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			w.Println(string(r))
			w.Println(err.Error())
			w.errorCount++
			return
		}

		var status struct {
			Status string `json:"status"`
		}

		json.Unmarshal(r, &status)
		if err != nil {
			w.Println(err.Error())
			w.errorCount++
			return
		}

		if status.Status != "success" {
			w.Println("Scalry Post Status:", status.Status)
			w.errorCount++
			return
		}
	}

	w.errorCount = 0
	w.buffer.Reset()
}

// Update periodically polls the current writer's buffer and uploads the logs to Scalyr.
// Interval should be provided in Milliseconds.
//
// Example usage:
// w := CreateScalyrWriter(os.Stdout, "https://www.scalyr.com/api/uploadLogs?host=ExampleService&logfile=AccessLog&token=ExampleToken")
// go w.Update(2000) // Update every 2 seconds.
// Client timeout is set to 2 seconds, it's not recommended that the update interval be less than 2s.
func (w *ScalyrWriter) Update(interval int) {
	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)

	for range ticker.C {
		w.UpdateNow()
	}
}

// Println writes a log message out to the defined TEE writer.
// Anything sent to Println will *NOT* be sent to scalyr.
func (w *ScalyrWriter) Println(msg ...string) {
	if len(msg) < 1 {
		return
	}

	m := strings.Join(msg, " ")
	w.tee.Write([]byte(m))
	w.tee.Write([]byte("\n"))
}
