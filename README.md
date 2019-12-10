# urlfuzz - WIP

urlfuzz is a highly concurrent AWS S3 URL fuzzer for time critical use cases. A distributed, clustered version that can can leverage as many as necessary workers nodes with load balancing is included. A AWS lambda version is included as proof of concept.

Benchmarked with 1.9 million operations per seconds running as cluster with master and worker nodes on 20 Amazon EC2 M6g instances.
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
