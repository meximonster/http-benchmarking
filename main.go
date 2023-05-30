package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.internal.upstreamsystems.com/george.petropoulos/http-benchmarking/configuration"
	"gitlab.internal.upstreamsystems.com/george.petropoulos/http-benchmarking/req"
	"gitlab.internal.upstreamsystems.com/george.petropoulos/http-benchmarking/srv"
)

var (
	c       *configuration.Config
	jbody   []byte
	freq    time.Duration
	buckets []float64
)

func init() {
	err := configuration.Load()
	if err != nil {
		log.Fatal("error loading configuration: ", err)
	}
	c = configuration.Read()
	jbody, freq, buckets, err = c.Validate()
	if err != nil {
		log.Fatal("error validating configuration: ", err)
	}
	req.RegisterMetrics(buckets)
}

func main() {

	go func() {
		err := srv.Run(c.Ip, c.Port)
		if err != nil {
			log.Println("http server returned error: ", err)
		}
	}()

	req.InitClient(time.Duration(c.ClientTimeout) * time.Second)

	log.Printf("Ready to fire %d requests from %d threads every %s to %s\n", c.Requests, c.Threads, c.Frequency, c.Endpoint)
	go req.Schedule(c.Threads, c.Requests, c.Endpoint, c.Method, jbody, freq, c.Verbose, c.UUIDParam)

	// Create signaling for process termination.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	srv.Close()
	req.Done()

}
