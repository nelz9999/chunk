# `chunk`
The chunk command line utility enables a user to add delays between streaming of chunks of bytes from the source stream

This may be helpful to people who are trying to simulate the output of live streaming applications, like audio or video streaming.

# Help output
```
$ chunk -h
The chunk command line utility enables a user to add delays between
streaming of chunks of bytes from the source stream

Usage:
  chunk [flags]

Flags:
  -d, --debug          send debugging output to stderr
  -h, --help           help for chunk
  -i, --input string   specify source file, otherwise defaults to stdin
  -l, --low-size int   set to a non-zero value less than the max-size to send random variable sized chunks of bytes
  -s, --max-size int   set the maximum chunk size, in bytes, to send (default 16)
  -w, --max-wait int   set the period, in milliseconds, to wait between chunk delivery (default 100)
  -m, --min-wait int   set to a non-zero value less than the max-wait to wait random variable periods between chunks
```

# Example usage
First, let's create a file filled with 1024 bytes of arbitrary data.
```
$ head -c 1024 /dev/urandom > example.bin
```

Copying the contents of the file to another file is (too) quick.
```
$ time cat example.bin > example.zero

real  0m0.009s
user  0m0.003s
sys   0m0.003s
```

We'll use `chunk`, ingesting via `stdin`, to make it slower.
```
$ time cat example.bin | chunk > example.one

real  0m6.659s
user  0m0.010s
sys   0m0.016s
```

We can also use `chunk` directly via file input.
```
$ time chunk -i example.bin > example.two

real  0m6.615s
user  0m0.007s
sys   0m0.014s
```

Just to be sure, we can check that the contents are all the same.
```
$ shasum -a 256 example.*
f062c239c5ee21d2e55812ed7de3676a08e2de88e91a304716a6838b52b37e4f  example.bin
f062c239c5ee21d2e55812ed7de3676a08e2de88e91a304716a6838b52b37e4f  example.one
f062c239c5ee21d2e55812ed7de3676a08e2de88e91a304716a6838b52b37e4f  example.two
f062c239c5ee21d2e55812ed7de3676a08e2de88e91a304716a6838b52b37e4f  example.zero
```
