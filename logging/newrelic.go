package logging

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	newrelic "github.com/newrelic/go-agent"
)

// NewRelicLogger creates a New Relic logger handler using the provided API credentials.
// Currently NewRelicLogger relies on functionality provided by gorilla/mux to retrieve the
// path template.
func NewRelicLogger(appName, license string) func(http.Handler) http.Handler {

	config := newrelic.NewConfig(appName, license)
	app, err := newrelic.NewApplication(config)
	if err != nil {
		log.Fatal(err)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if app == nil {
				next.ServeHTTP(w, r)
				return
			}

			route := mux.CurrentRoute(r)
			if route == nil {
				next.ServeHTTP(w, r)
				return
			}

			template, err := route.GetPathTemplate()
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			txn := app.StartTransaction(template, w, r)
			defer txn.End()

			r = newrelic.RequestWithTransactionContext(r, txn)
			next.ServeHTTP(txn, r)
		})
	}
}
