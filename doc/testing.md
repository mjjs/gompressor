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
E.coli | 4,638,690 bytes | The complete genome of the E. Coli bacterium. Highly repetitive data.
world192.txt | 2,473,400 bytes | The CIA world fact book.
Randomdata | 2000000 bytes | Random bytes for testing patternless compression.

The random data has been generated with the command `head -c 2000000 </dev/urandom >Randomdata`.

### Results

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
The rest of the tests will be using the largest dictionary size.

Now, lets compare how the different files affect the algorithm's compression ratio. The `E.coli` file obviously
has lots of repetition, whereas a file filled with random data contains next to no repetition.  The CIA
world fact book is "regular" english text, so it will stand somewhere between the other two files.

filename     | original size (bytes) | compressed size (bytes)   | compression ratio (% of original) | compress time (µs)
-------------|-----------------------|---------------------------|-----------------------------------|------------------
E.coli       | 4,638,690             | 1,342,580                 | 28.94                             | 4,215,383
world192.txt | 2,473,400             | 1,076,126                 | 43.51                             | 8,248,130
Randomdata   | 2,000,000             | 2,934,700                 | 146.73                            | 7,250,137,6

As might be expected, the data with lots of repetition gets compressed the most, whereas the randomly
generated data is actually not compressed at all. Instead it turns out to get larger than the original.
The CIA world fact book might be a good indicator of what kind of compression ratio we can expect when
compressing text. The randomly generated data also took a much longer time to "compress" than the data
with repetition.

Finally, let's see how the input file size affects the compression time and ratio for the LZW algorithm.

filename       | original size (bytes)  | compressed size (bytes) | compression ratio (% of original) | compress time (µs)
---------------|------------------------|-------------------------|-----------------------------------|-------------------
world192.txt   | 1,024                  | 1,208                   | 117.97                            | 8,02
world192.txt   | 2,048                  | 2,206                   | 107.71                            | 1,465
world192.txt   | 4,096                  | 3,942                   | 96.24                             | 2,943
world192.txt   | 8,192                  | 6,848                   | 83.59                             | 5,373
world192.txt   | 1,638,4                | 1,254,6                 | 76.57                             | 1,193,9
world192.txt   | 3,276,8                | 2,251,4                 | 68.71                             | 2,224,9
world192.txt   | 6,553,6                | 3,927,6                 | 59.93                             | 5,224,1
world192.txt   | 1,310,72               | 6,852,4                 | 52.28                             | 1,068,09
world192.txt   | 2,621,44               | 1,189,56                | 45.38                             | 2,285,32
world192.txt   | 5,242,88               | 2,345,70                | 44.74                             | 4,903,35
world192.txt   | 1,048,576              | 4,664,96                | 44.49                             | 9,660,38
world192.txt   | 2,097,152              | 9,125,52                | 43.51                             | 2,073,592
world192.txt   | 2,473,400              | 1,076,126               | 43.51                             | 2,377,361

From the results, we can observe that the LZW algorithm actually works better for larger inputs.

#### Huffman
As there are no parameters to change in the Huffman compressing algorithm (at least my implementation),
I started by seeing how the different kinds of data affect the compression ratio of the algorithm.

filename       | original size (bytes)  | compressed size (bytes) | compression ratio (% of original) | compress time (µs)
---------------|------------------------|-------------------------|-----------------------------------|-------------------
E.coli         | 4,638,690              | 1,159,685               | 25.00                             | 2,932,067
world192.txt   | 2,473,400              | 1,558,877               | 63.03                             | 2,495,041
Randomdata     | 2,000,000              | 2,000,768               | 100.04                            | 3,058,153

We can see that the results are similar to the results we got with Lempel-Ziv-Welch. Data with high repetition gets compressed
more than data with less repetition. Interestingly, we see that the randomly generated data gets compressed to around the same
size as the original size, so it is much more efficient than LZW (even though the result is still bigger than the input).
We can also see that the compression remain consistent across different types of input. Compression also takes much less
time compared to LZW. My implementation of the LZW is most likely at fault here, though.

Next, let's inspect how the size of the original data affects the compression time and ratio of the Huffman algorithm.

filename       | original size (bytes)  | compressed size (bytes) | compression ratio (% of original) | compress time (µs)
---------------|------------------------|-------------------------|-----------------------------------|-------------------
world192.txt   |      1,024             | 819                     | 79.98                             | 380
world192.txt   |      2,048             | 1,515                   | 73.97                             | 645
world192.txt   |      4,096             | 2,821                   | 68.87                             | 1,326
world192.txt   |      8,192             | 5,384                   | 65.72                             | 2,337
world192.txt   |      1,638,4           | 1,062,2                 | 64.83                             | 4,556
world192.txt   |      3,276,8           | 2,092,7                 | 63.86                             | 9,426
world192.txt   |      6,553,6           | 4,144,3                 | 63.24                             | 2,632,6
world192.txt   |      1,310,72          | 8,209,6                 | 62.63                             | 4,380,4
world192.txt   |      2,621,44          | 1,651,55                | 63.00                             | 8,988,0
world192.txt   |      5,242,88          | 3,306,15                | 63.06                             | 1,799,39
world192.txt   |      1,048,576         | 6,620,16                | 63.13                             | 3,429,33
world192.txt   |      2,097,152         | 1,323,510               | 63.11                             | 6,995,93
world192.txt   |      2,473,400         | 1,558,877               | 63.03                             | 7,396,87

We can see that the compression ratio stays quite conistent across different input sixes. The compress time, however
grows quite drastically.
