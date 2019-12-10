package main

import (
	"os"
	"testing"

	"github.com/paulbellamy/ratecounter"
)

func Test_isLetter(t *testing.T) {
	type args struct {
		c rune
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isLetter(tt.args.c); got != tt.want {
				t.Errorf("isLetter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_httpRequest(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := httpRequest(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("httpRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("httpRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_producer(t *testing.T) {
	type args struct {
		jobs             chan<- *Job
		stopChan         chan struct{}
		baseURL          string
		maxNumbers       int
		availableLetters string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			producer(tt.args.jobs, tt.args.stopChan, tt.args.baseURL, tt.args.maxNumbers, tt.args.availableLetters)
		})
	}
}

func Test_worker(t *testing.T) {
	type args struct {
		jobs    <-chan *Job
		results chan<- *Job
		counter *ratecounter.RateCounter
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			worker(tt.args.jobs, tt.args.results, tt.args.counter)
		})
	}
}

func Test_signalStop(t *testing.T) {
	type args struct {
		signalCh chan os.Signal
		stopChan chan struct{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signalStop(tt.args.signalCh, tt.args.stopChan)
		})
	}
}

func Test_rateDisplay(t *testing.T) {
	type args struct {
		counter *ratecounter.RateCounter
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rateDisplay(tt.args.counter)
		})
	}
}

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}
