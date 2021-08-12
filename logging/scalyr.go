package logging

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/btm6084/gojson"
	"github.com/spf13/cast"
)

var (
	// ScalyrClient is the http.Client used to communicate with scalyr.
	ScalyrClient = &http.Client{
		Timeout: 2 * time.Second,
	}

	scalyrBufferLimit = 50000000
)

// ScalyrWriter will buffer all log writes for logging to scalyr, then
// tee all calls to the supplied io.Writer
type ScalyrWriter struct {
	// tee provides an io.Writer that receives all log entries, along with the entries written to scalyr.
	// popular options include os.Stdout and ioutil.Discard (/dev/null)
	tee            io.Writer
	buffer         bytes.Buffer
	url            string
	lock           sync.Mutex
	lastReset      time.Time
	updateInterval time.Duration

	scalyrBuffer [][]byte
	scalyrLock   sync.Mutex
}

func (w *ScalyrWriter) scalyrBufferLen() int {
	sum := 0

	for i := 0; i < len(w.scalyrBuffer); i++ {
		sum += len(w.scalyrBuffer[i])
	}

	return sum
}

func (w *ScalyrWriter) Write(p []byte) (int, error) {
	// If we can't upload to Scalyr and flush the buffer, we stop accumulating
	// so that we don't overrun memory by infinitely buffering.
	if w.scalyrBufferLen() <= scalyrBufferLimit {
		w.lock.Lock()
		w.buffer.Write(p)

		if w.buffer.Len() > 5000000 {
			buf := make([]byte, w.buffer.Len())
			copy(buf, w.buffer.Bytes())

			w.scalyrLock.Lock()
			w.scalyrBuffer = append(w.scalyrBuffer, buf)
			w.scalyrLock.Unlock()

			w.buffer.Reset()
			w.lastReset = time.Now()
		}

		w.lock.Unlock()
	}
	n, err := w.tee.Write(p)

	return n, err
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

	return &ScalyrWriter{tee: tee, url: url, updateInterval: 0, lastReset: time.Now()}
}

// UpdateNow immediately uploads the collected logs to Scalyr.
// Interval should be provided in Milliseconds.
//
// Example usage:
// w := CreateScalyrWriter(os.Stdout, "https://www.scalyr.com/api/uploadLogs?host=ExampleService&logfile=AccessLog&token=ExampleToken")
// defer w.UpdateNow() // Update after leaving the current function.
func (w *ScalyrWriter) UpdateNow(flushBuffer bool) {
	w.scalyrLock.Lock()
	defer w.scalyrLock.Unlock()

	old := time.Since(w.lastReset) >= w.updateInterval
	empty := w.buffer.Len() == 0
	full := w.scalyrBufferLen() > scalyrBufferLimit

	if flushBuffer || (old && !empty && !full) {
		w.lock.Lock()
		buf := make([]byte, w.buffer.Len())
		copy(buf, w.buffer.Bytes())

		w.scalyrBuffer = append(w.scalyrBuffer, buf)
		w.buffer.Reset()
		w.lastReset = time.Now()
		w.lock.Unlock()
	}

	for len(w.scalyrBuffer) > 0 {
		buf := w.scalyrBuffer[0]

		start := time.Now()
		resp, err := ScalyrClient.Post(w.url, "text/plain", bytes.NewBuffer(buf))
		if err != nil {
			w.Println(err.Error())
			return
		}
		w.Println("Scalyr Post Time:", time.Since(start).String(), "Post Size:", cast.ToString(len(buf)), "Scalyr Post Url:", w.url)

		r, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			w.Println(string(r))
			w.Println(err.Error())
			return
		}

		w.Println("Scalyr Response:", string(r))

		status, err := gojson.ExtractString(r, "status")
		if status != "success" {
			errmsg := "non-success status from scalyr"
			if err != nil {
				errmsg = err.Error()
			}

			w.Println("Scalry Post Status:", status, errmsg)
			w.Println(string(r))
			return
		}

		w.scalyrBuffer = w.scalyrBuffer[1:]
	}
}

// Update periodically polls the current writer's buffer and uploads the logs to Scalyr.
// Interval should be provided in Milliseconds.
//
// Example usage:
// w := CreateScalyrWriter(os.Stdout, "https://www.scalyr.com/api/uploadLogs?host=ExampleService&logfile=AccessLog&token=ExampleToken")
// go w.Update(2000) // Update every 2 seconds.
// Client timeout is set to 2 seconds, it's not recommended that the update interval be less than 2s.
func (w *ScalyrWriter) Update(interval int) {
	w.updateInterval = time.Duration(interval) * time.Millisecond
	// Update no more often than 1 time per second.
	if w.updateInterval < 1*time.Second {
		w.updateInterval = 1 * time.Second
	}

	ticker := time.NewTicker(w.updateInterval)

	for range ticker.C {
		w.UpdateNow(false)
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
