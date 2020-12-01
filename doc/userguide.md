# User guide

The executable for the program can be found over at the releases section of the
repository in Github: https://github.com/mjjs/gompressor/releases/latest.

## Running the program
The program has two ways to use it. When starting the program with the `-tui` flag,
the application starts in a text-based user interface. The interface is very simple,
and allows for compression and decompression of files found in the current directory
and any subdirectory.

To run the application without any user interface, you need to supply a few command
line flags. The command line flags can be printed by supplying the `-help` flag
to the application. Here are some example use cases:

```bash
# Compressing a file using the Lempel-Ziv-Welch algorithm
./gompressor -lzw -compress -in=/path/to/input/file -out=/path/to/save/compressed/file/into
```

```bash
# Decompressing a file compressed using the Huffman algorithm
./gompressor -huffman -decompress -in=/path/to/compressed/file -out=/path/to/save/decompressed/file/into
```

## Inputs
As the compression algorithms work on bytes, in theory, the program can compress
any file that is given to it. In practice, however, I found that compressing larger
than around 5MB non-text files takes quite a long time. The testdata folder contains
text files for testing the algorithms.
