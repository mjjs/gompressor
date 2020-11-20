# Testing documentation - WORK IN PROGRESS
## Unit testing
The unit tests are written using Golang's own testing functionality. All the data
structures and core algorithms have tests to verify their functionality, and
the test coverage for these packages is close to 100%.

If Go is installed, the tests can be ran from the root of the project by running
`go test ./...` or `make test`. The current coverage is automatically updated to
[Coveralls](https://coveralls.io/github/mjjs/gompressor?branch=master) after each merge
to master. There is also a local [coverage.html](./coverage.html) file in the repository which
I try to keep up to date so it can be viewed offline.

## Performance testing
The performance testing is done by `compressiontester`, which can be found in a
separate folder in the root of the project. The tester uses three different files
and runs compression and decompression using the LZW and Huffman algorithms and
writes the output into csv files.

The test files are:
filename | size | description
---------|------|------------
E.coli | 4,638,690 bytes | The complete genome of the E. Coli bacterium
world192.txt | 2,473,400 bytes | The CIA world fact book
Randomdata | 2000000 bytes | Random bytes for testing patternless compression

The random data has been generated with the command `head -c 2000000 </dev/urandom >Randomdata`.

### Results

#### Huffman

#### Lempel-Ziv-Welch
First I tested how the dictionary size of the LZW algorithm affects the algorithm.
Here in the table we can see the results on compressing the E.coli text file
with LZW using different dictionary sizes.

dictionary size | original size (bytes) | compressed size (bytes) | compression ratio (% of original) | compress time (µs)
----------------|-----------------------|-------------------------|-----------------------------------|--------------
512             | 4,638,690             | 3,184,280               | 68.65                             | 4,160,141
1023            | 4,638,690             | 2,523,088               | 54.39                             | 3,414,906
4095            | 4,638,690             | 1,917,598               | 41.34                             | 3,336,908
32767           | 4,638,690             | 1,448,820               | 31.23                             | 3,810,436
65535           | 4,638,690             | 1,342,580               | 28.94                             | 4,215,383

As we can see, the compression ratio improves noticeably when growing the dictionary size.
