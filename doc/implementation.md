# Implementation documentation

## Project structure
All the code in this project are separated in their own packages, which are named
after what they implement.

### Data structures

#### dictionary
Dictionary implements a hash table which is used by both, the LZW and Huffman
compression algorithms.

#### linkedlist
A linked list which is used in the dictionary as a fallback for hash collisions.

#### priorityqueue
A priority queue, which is implemented using a minimum binary heap. The priority queue
is used by the Huffman algorithm when creating the prefix tree used in data compression.

#### vector
Vector is a dynamic array. It is implemented as an array list, and provides O(1) access
to the elements.

### Compression algorithms

#### lzw
Implements the Lempel-Ziv-Welch lossless data compression algorithm. The compression
takes in a vector of bytes and outputs a vector of 16-bit integer codes. The decompression
algorithm takes in the vector of said 16-bit codes and outputs a vector of bytes, which
is identical to the vector initially given to the compress algorithm.

The algorithm works by building a dictionary of sequences of bytes encountered earlier
in the input vector. These sequences get their 16-bit codes from the current length of
the dictionary when the sequence was first discovered.

#### huffman
Implements the Huffman coding lossless data compression algoritihm. Input to the
compression algorithm is a vector of bytes, and the output is a vector of compressed
bytes. The decompression algorithm also takes in a vector of bytes and decompresses
it into the original vector of bytes.

Compression starts by building a prefix tree of the bytes found in the input vector
so that the more frequent the byte is in the input, the higher the byte will be in
the prefix tree. This ensures that often appearing bytes will be compressed into
fewer bits. After this, a Huffman code is calculated for each byte in the input
by walking the tree and appending 0 or 1 to the code depending on which leaf the
algorithm follows. These codes are then used to compress the original bytes. Finally,
the prefix tree along with the huffman codes are written to the output vector.

#### File I/O
to be written

### Time complexities

#### Lempel-Ziv-Welch
The time complexity of the algorithm is O(n), as the input is only iterated once,
and the codes are calculated in effectively O(1) time.

#### Huffman
The building of the tree is done in O(n log n) time, and going through each byte
in the input to find their Huffman code takes O(n) time. The complexity of the
whole algorithm is O(n log n).
