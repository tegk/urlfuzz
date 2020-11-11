# urlfuzz

urlfuzz is a highly concurrent AWS S3 URL fuzzer for time-critical use cases. A distributed, clustered version, which can use as many worker nodes as necessary with load balancing, is included. An AWS lambda version is included as proof of concept.

Benchmarking with 1.9 million operations per second running as a cluster with master and worker nodes on 20 Amazon EC2 M6g instances.
## Installation

Build it from source:

```bash
git clone https://github.com/tegk/urlfuzz
cd urlfuzz
go build
```

## Usage

```go
        defaultBaseUrl := "https://test-assets.s3.amazonaws.com/test/test/20190619/20190619-TEST-test%s%s%s%03d.png"
	workerRoutines := flag.Int("threads", 5000, "")
	maxNumbers := flag.Int("maxNumbers", 1000, "")
	maxJobs := flag.Int("maxjobs", 1000000, "")
	availableLetters := flag.String("availableLetters", "abcdefghijklmnopqrstuvwxyz", "")
	baseUrl := flag.String("baseURL", defaultBaseUrl, "")
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)
