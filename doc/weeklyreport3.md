# Week 3

## What was done
During this week, I have successfully implemented the Huffman coding algorithm. I did not have too much time to work
on the project due to other projects taking a lion's share of my time. This should, however, not be a problem next week.

## What I learned
I feel like I understand the Huffman coding algorithm fairly well now.
After reading about it and implementing it, turns out it is quite simple in practice.

## What's next
As I have now implemented both of the intended compression algorithms, the next step will be to refactor the code to
be more readable and increase the test coverage as much as possible.

I will also start working on ways to actually compare the compression algorithms. One idea is to parametrize the
LZW algorithm's dictionary size to see how much it affects the compression ratio.

## What has been problematic
Not so much of a problem, but my first implementation of the Huffman compression algorithm ended up actually creating
a file that was ten times as large as the original file. This was actually a known issue during implementation time,
because I wrote each Huffman code as separate bytes instead of bits. This has been fixed since, but the fix manifested
itself as some not-so-pretty bitwise manipulation code, which is hard to read currently.

## Time used
time | activity |
-----|----------|
  1h | Studying the Huffman coding algorithm further |
  3h | Implementing the Huffman coding algorithm     |
  1h | Improving Huffman compression ratio           |
