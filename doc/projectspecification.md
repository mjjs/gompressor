# Project specification

The goal of this project is to implement at least two different compression algorithms and compare how well they perform.

The algorithms that will be implemented are the Huffman coding algorithm and the Lempel-Ziv-Welch (LZW) algorithm. LZW was chosen over LZ77 and LZ78 because it should be more efficient, yet still fairly straightforward to implement.

## Administrative information
Degree programme: bachelor's degree in computer science

Documentation language: English

Programming language: Golang

## Data structures
### LZW
The algorithm will use a hashtable as a dictionary to map data sequences into codes. A dynamically sized array will also be used as an auxiliary data structure.

### Huffman
The Huffman algorithm will will utilize a priority queue implemented as a min heap.

## Time complexity
### LZW
The time complexity should be O(n) since the input string is only iterated through once according to my initial research of the algorithm.

### Huffman
The tree-structure of the Huffman algorithm requires O(n log n) time.

## Input and output
The input of the program will be a file, and the output will be a compressed file. The output file should be noticeably more smaller than the input file.

## Sources
* https://en.wikipedia.org/wiki/Lempel%E2%80%93Ziv%E2%80%93Welch
* https://www2.cs.duke.edu/csed/curious/compression/lzw.html
* https://ocw.mit.edu/courses/electrical-engineering-and-computer-science/6-046j-design-and-analysis-of-algorithms-spring-2012/lecture-notes/MIT6_046JS12_lec19.pdf
* https://en.wikipedia.org/wiki/Huffman_coding
* http://warp.povusers.org/EfficientLZW/index.html
