package main

import (
	"flag"
	"fmt"
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

func isLetter(c rune) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
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

func producer(jobs chan<- *Job, stopChan chan struct{}, baseURL string, maxNumbers int, availableLetters string) {
	alphabetLen := len(availableLetters)
	//var links []string
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
					newURL := fmt.Sprintf(baseURL, string(firstLetter), string(secondLetter), string(thirdLetter), number)
					jobs <- &Job{url: newURL}
				}
			}
		}
	}
	close(jobs)
}

func worker(jobs <-chan *Job, results chan<- *Job, counter *ratecounter.RateCounter) {
	defer wg.Done()
	for {
		j, more := <-jobs
		if more {
			_, err := httpRequest(j.url)
			counter.Incr(1)
			if err == nil {
				log.Println("found valid url:", j.url)
				os.Exit(2)
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

func main() {
	//"https://usg-assets.s3.amazonaws.com/assets/upbit-sg/201911/20191121-usg_%s%s%s%03d.png"
	//defaultBaseUrl := "http://localhost:2027/assets/upbit-sg/201911/20191121-usg_%s%s%s%03d.jpg"

	defaultBaseUrl := "https://pbg-assets.s3.amazonaws.com/editorial/pbo/20190619/20190619-PBO-Update_%s%s%s%03d.png"

	workerRoutines := flag.Int("threads", 5000, "")
	maxNumbers := flag.Int("maxNumbers", 1000, "")

	maxJobs := flag.Int("maxjobs", 1000000, "")
	availableLetters := flag.String("availableLetters", "abcdefghijklmnopqrstuvwxyz", "")
	baseUrl := flag.String("baseURL", defaultBaseUrl, "")
	flag.Parse()

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

	go producer(jobs, stopChan, *baseUrl, *maxNumbers, *availableLetters)
	go signalStop(signalCh, stopChan)
	go rateDisplay(counter)
	wg.Wait()
	close(results)
	<-done
}
