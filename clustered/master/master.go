package main

import (
	"flag"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"os/signal"
	"sync"
)

var wg sync.WaitGroup

type Job struct {
	url []byte
}

func isLetter(c rune) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}

func producer(jobs chan<- *Job, stopChan chan struct{}, baseURL string, maxNumbers int, availableLetters string) {
	alphabetLen := len(availableLetters)
	for _, firstLetter := range availableLetters {
		if !isLetter(firstLetter) {
			log.Fatalf("a number or non alphabetic character was entered in letters flag: %v", firstLetter)
		}
		secondLetter := 'a'
		thirdLetter := 'a'
		for second := 0; second < alphabetLen; second++ {
			secondLetter = rune(availableLetters[second])
			for third := 0; third < alphabetLen; third++ {
				thirdLetter = rune(availableLetters[third])
				for number := 0; number < maxNumbers; number++ {
					newUrl := fmt.Sprintf(baseURL, string(firstLetter), string(secondLetter), string(thirdLetter), number)
					jobs <- &Job{url: []byte(newUrl)}
				}
			}
		}
	}
	close(jobs)
}

func printMsg(m *nats.Msg, i int) {
	log.Printf("[#%d] Received on [%s]: '%s'", i, m.Subject, string(m.Data))
}

func resultSubscriber(done chan<- bool, nc *nats.Conn) {
	nc.Subscribe("valid", func(msg *nats.Msg) {})
	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}
}

func publisher(jobs <-chan *Job, results chan<- *Job, nc *nats.Conn) {
	defer wg.Done()
	for {
		j, more := <-jobs
		if more {
			err := nc.Publish("test", j.url)
			fmt.Println(string(j.url))
			if err != nil {
				log.Println(err)
			}
		} else {
			nc.Flush()
			return
		}
	}
}

func signalStop(signalCh chan os.Signal, stopChan chan struct{}) {
	<-signalCh
	close(stopChan)
}

func main() {
	//"https://usg-assets.s3.amazonaws.com/assets/upbit-sg/201911/20191121-usg_%s%s%s%03d.png"
	//defaultBaseUrl := "http://localhost:2027/assets/upbit-sg/201911/20191121-usg_%s%s%s%03d.jpg"

	var (
		defaultBaseUrl   = "https://pbg-assets.s3.amazonaws.com/editorial/pbo/20190619/20190619-PBO-Update_%s%s%s%03d.png"
		workerRoutines   = flag.Int("threads", 1, "")
		maxJobs          = flag.Int("maxjobs", 500000, "")
		availableLetters = flag.String("availableLetters", "abcdefghijklmnopqrstuvwxyz", "")
		baseUrl          = flag.String("baseURL", defaultBaseUrl, "")
		maxNumbers       = flag.Int("maxNumbers", 1000, "")
		urls             = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")
	)

	flag.Parse()

	//nats Connect Options
	opts := []nats.Option{nats.Name("master")}

	// Connect to NATS
	nc, err := nats.Connect(*urls, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	var jobs = make(chan *Job, *maxJobs)
	var results = make(chan *Job)
	var done = make(chan bool, 1)

	stopChan := make(chan struct{})

	signalCh := make(chan os.Signal, 2)
	signal.Notify(signalCh, os.Interrupt)

	for w := 1; w <= *workerRoutines; w++ {
		wg.Add(1)
		go publisher(jobs, results, nc)
	}

	go producer(jobs, stopChan, *baseUrl, *maxNumbers, *availableLetters)
	go resultSubscriber(done, nc)
	go signalStop(signalCh, stopChan)
	wg.Wait()
	close(results)
	<-done
}
