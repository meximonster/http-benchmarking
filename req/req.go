package req

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"gitlab.internal.upstreamsystems.com/george.petropoulos/http-benchmarking/metrics"
)

var (
	// reqTotal counts total exports based on the status code.
	reqTotal *prometheus.CounterVec
	// reqLatency initiates a histogram to count the latency of the requests.
	reqLatency *prometheus.HistogramVec

	sum prometheus.Summary
)

func RegisterMetrics(buckets []float64) {
	reqTotal = metrics.NewCounterVec("total_requests", "The total number of requests", []string{"status", "path"})
	reqLatency = metrics.NewHistogramVec("request_latency_ms", "The latency of the requests", buckets, "code", "path")
	sum = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name:       "api_request_durations",
			Help:       "total request latency",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
	)
	prometheus.MustRegister(reqTotal)
	prometheus.MustRegister(reqLatency)
	prometheus.MustRegister(sum)
}

func Make(wg *sync.WaitGroup, ch chan *http.Request, requests int, url string, method string, jbody []byte, param bool) {
	defer wg.Done()
	defer close(ch)
	for j := 0; j < requests; j++ {
		req, err := http.NewRequest(method, url, bytes.NewBuffer(jbody))
		if err != nil {
			log.Println("error creating request:", err)
			continue
		}
		if param {
			uuid := uuid.New()
			params := req.URL.Query()
			params.Add("uuid", uuid.String())
			req.URL.RawQuery = params.Encode()
		}
		ch <- req
	}
}

func Do(wg *sync.WaitGroup, ch chan *http.Request, client http.Client, threadNumber int, verbose bool, param bool) {
	defer wg.Done()
	for req := range ch {
		req.Header.Set("Content-Type", "application/json")
		start := time.Now()
		timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
			us := v * 1000 // make milliseconds
			sum.Observe(us)
		}))
		resp, err := client.Do(req)
		if err != nil {
			reqTotal.WithLabelValues("fail", req.URL.Path).Inc()
			log.Println("error performing request:", err)
			continue
		}
		timer.ObserveDuration()
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
			reqLatency.WithLabelValues("error", req.URL.Path).Observe(float64(time.Since(start).Milliseconds()))
			reqTotal.WithLabelValues("error", req.URL.Path).Inc()
		}
		reqLatency.WithLabelValues("success", req.URL.Path).Observe(float64(time.Since(start).Milliseconds()))
		reqTotal.WithLabelValues("success", req.URL.Path).Inc()
		if verbose {
			s := fmt.Sprintf("thread #%d latency %vms", threadNumber, float64(time.Since(start).Milliseconds()))
			if param {
				id := req.URL.Query().Get("uuid")
				s += " id " + id
			}
			log.Println(s)
		}
		resp.Body.Close()
	}
}
