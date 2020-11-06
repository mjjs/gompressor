# Week 2

## What was done
During this week, I have managed to replace the standard library data structures used in the LZW encoding/decoding with ones implemented by me.
I have also implemented a priority queue to be used with Huffman coding. The initial work for the Huffman coding algorithm has started, but that code is not merged into master yet.

## What I learned
I have gained some insight into how the Huffman coding algorithm works. I can't say I know 100% how it works, but I have a basic understanding.
I also feel like my understanding of the LZW encoding when I was replacing the standard library data structures with my own data structures.

## What's next
I will continue working on the Huffman compression and decompression algorithms. As I have most of the required data structures already implemented, I think it should be fairly simple to implement it using those instead of using standard library data structures.

## What has been problematic
Nothing has been too problematic so far. The only thing that caused me some trouble was implementing the priority queue.
This was mainly because I had forgotten the intricacies of the heap data structure, so I had to spend some time studying it before I got it to work properly.

## Questions
* What kind of UI is the "minimum" for the project?

Does running the project via command line flags suffice, or should I implement some kind of text-based interface?

## Time used
time | activity |
-----|----------|
  0.25h | Setting up Travis and Coveralls for the project     |
  3h    | Implementing data structures for the LZW algorithms |
  2h    | Implementing the priority queue                     |
