package srv

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gitlab.internal.upstreamsystems.com/george.petropoulos/http-benchmarking/metrics"
)

var srv *http.Server

func Run(ip string, port string) error {

	r := mux.NewRouter()
	srv = &http.Server{
		Handler:      r,
		Addr:         ip + ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	r.Handle("/metrics", metrics.Handler()).Methods("GET")

	if err := srv.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func Close() {
	srv.Close()
}
