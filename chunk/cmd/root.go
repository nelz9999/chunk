// Copyright © 2017 Nelz
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/nelz9999/chunk/stream"
	"github.com/spf13/cobra"
)

var debug bool
var maxSize int
var lowSize int
var maxWait int
var lowWait int
var input string

// RootCmd represents the base "chunk" command
var RootCmd = &cobra.Command{
	Use:   "chunk",
	Short: "Add chunkiness and delay to a stream",
	Long: `The chunk command line utility enables a user to add delays between
streaming of chunks of bytes from the source stream`,
	RunE: func(cmd *cobra.Command, args []string) error {
		in, err := buildReader()
		if err != nil {
			return err
		}
		defer in.Close()

		log := buildLog()

		sizer, buf, err := buildSizer()
		if err != nil {
			return err
		}

		waiter, err := buildWaiter()
		if err != nil {
			return err
		}

		src := stream.New(in, sizer, waiter, log)

		_, err = io.CopyBuffer(os.Stdout, src, buf)
		return err
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "send debugging output to stderr")
	RootCmd.Flags().IntVarP(&maxSize, "max-size", "s", 16, "set the maximum chunk size, in bytes, to send")
	RootCmd.Flags().IntVarP(&lowSize, "low-size", "l", 0, "set to a non-zero value less than the max-size to send random variable sized chunks of bytes")
	RootCmd.Flags().IntVarP(&maxWait, "max-wait", "w", 100, "set the period, in milliseconds, to wait between chunk delivery")
	RootCmd.Flags().IntVarP(&lowWait, "min-wait", "m", 0, "set to a non-zero value less than the max-wait to wait random variable periods between chunks")
	RootCmd.Flags().StringVarP(&input, "input", "i", "", "specify source file, otherwise defaults to stdin")
}

// buildReader prepares the input stream
func buildReader() (io.ReadCloser, error) {
	result := os.Stdin

	if input != "" {
		var err error
		result, err = os.Open(input)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

// buildLog prepares the writer that receives the debugging output
func buildLog() io.Writer {
	result := ioutil.Discard
	if debug {
		result = os.Stderr
	}
	return result
}

// buildSizer prepares the the object that dictates how large the chunks
// are, and the buffer that gets used for the io.CopyBuffer
func buildSizer() (stream.Inter, []byte, error) {
	buf := make([]byte, maxSize)
	sizer := stream.InterFunc(func() int {
		return maxSize
	})

	if lowSize > maxSize {
		return nil, nil, fmt.Errorf("low-size > max-size: %d > %d", lowSize, maxSize)
	}

	if lowSize > 0 && lowSize < maxSize {
		breadth := maxSize - lowSize + 1
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		sizer = stream.InterFunc(func() int {
			return lowSize + r.Intn(breadth)
		})
	}

	return sizer, buf, nil
}

// buildWaiter prepares the object that dictates how long a period between
// chunks are emitted
func buildWaiter() (stream.Inter, error) {
	waiter := stream.InterFunc(func() int {
		return maxWait
	})

	if lowWait > maxWait {
		return nil, fmt.Errorf("min-wait > max-wait: %d > %d", lowWait, maxWait)
	}

	if lowWait > 0 && lowWait < maxWait {
		breadth := maxWait - lowWait + 1
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		waiter = stream.InterFunc(func() int {
			return lowWait + r.Intn(breadth)
		})
	}

	return waiter, nil
}
