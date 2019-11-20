package main

import (
	"flag"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/paulbellamy/ratecounter"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

var wg sync.WaitGroup

type Job struct {
	url string
}

func httpRequest(url string) (bool, error) {
	var netClient = &http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := netClient.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	_, err = io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		return false, err
	}
	if resp.StatusCode == http.StatusOK {
		return true, nil
	}
	return false, fmt.Errorf("wrong status")
}

func producer(jobs chan<- *Job, stopChan chan struct{}, nc *nats.Conn) {
	nc.QueueSubscribe("test", "queue", func(msg *nats.Msg) {
		jobs <- &Job{url: string(msg.Data)}
	})
	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}
}

func writer(results <-chan *Job, done chan<- bool, nc *nats.Conn, validSub string) {
	for {
		j, more := <-results
		if more {
			err := nc.Publish(validSub, []byte(j.url))
			if err != nil {
				log.Println(err)
			}
		} else {
			time.Sleep(time.Minute)
		}
	}
}

func worker(jobs <-chan *Job, results chan<- *Job, counter *ratecounter.RateCounter) {
	defer wg.Done()
	for {
		j, more := <-jobs
		if more {
			_, err := httpRequest(j.url)
			counter.Incr(1)
			if err == nil {
				results <- j
			}
		} else {
			return
		}
	}
}

func signalStop(signalCh chan os.Signal, stopChan chan struct{}) {
	<-signalCh
	close(stopChan)
}

func rateDisplay(counter *ratecounter.RateCounter) {
	for {
		time.Sleep(time.Second)
		fmt.Println(counter.Rate())
	}
}

func setupConnOptions(opts []nats.Option) []nats.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		log.Printf("Disconnected due to: %s, will attempt reconnects for %.0fm", err, totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Printf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		log.Fatalf("Exiting: %v", nc.LastError())
	}))
	return opts
}

func main() {
	var (
		workerRoutines = flag.Int("threads", 20, "")
		maxJobs        = flag.Int("maxjobs", 100000, "")
		urls           = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")
		validSub       = flag.String("validsubj", "valid", "subj to pub valid")
	)
	flag.Parse()

	opts := []nats.Option{nats.Name("sub1")}

	opts = setupConnOptions(opts)

	// Connect to NATS
	nc, err := nats.Connect(*urls, opts...)
	if err != nil {
		log.Fatal(err)
	}

	var jobs = make(chan *Job, *maxJobs)
	var results = make(chan *Job)
	var done = make(chan bool, 1)

	counter := ratecounter.NewRateCounter(1 * time.Second)

	stopChan := make(chan struct{})

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	for w := 1; w <= *workerRoutines; w++ {
		wg.Add(1)
		go worker(jobs, results, counter)
	}

	go producer(jobs, stopChan, nc)
	go writer(results, done, nc, *validSub)
	go signalStop(signalCh, stopChan)
	go rateDisplay(counter)
	wg.Wait()
	close(results)
	<-done
	nc.Drain()
}
