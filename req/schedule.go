package req

import (
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	done   = make(chan bool)
	client http.Client
)

func Done() {
	done <- true
}

func InitClient(seconds time.Duration) {
	client = http.Client{Timeout: seconds}
}

func Schedule(threads int, requests int, url string, method string, jbody []byte, frequency time.Duration, verbose bool, param bool) {
	ticker := time.NewTicker(frequency)
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			ch := make(chan *http.Request)
			if verbose {
				log.Println("Start processing...")
			}
			wg := sync.WaitGroup{}
			for i := 0; i < threads; i++ {
				wg.Add(1)
				go Do(&wg, ch, client, i+1, verbose, param)
			}
			wg.Add(1)
			go Make(&wg, ch, requests, url, method, jbody, param)
			wg.Wait()
			if verbose {
				log.Println("Finished processing")
			}
		}
	}
}
