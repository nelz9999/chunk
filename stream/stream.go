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

package stream

import (
	"fmt"
	"io"
	"time"
)

type Inter interface {
	Int() int
}

type InterFunc func() int

func (f InterFunc) Int() int {
	return f()
}

type ReaderFunc func([]byte) (int, error)

func (f ReaderFunc) Read(p []byte) (n int, err error) {
	return f(p)
}

// New creates a new io.Reader that spits out chunks of data after waiting
// a bit of time.
func New(r io.Reader, sizer Inter, waiter Inter, log io.Writer) io.Reader {
	return ReaderFunc(func(p []byte) (int, error) {
		wait := waiter.Int()
		max := sizer.Int()
		if len(p) < max {
			max = len(p)
		}
		n, err := r.Read(p[0:max])
		// TODO: Debug
		fmt.Fprintf(log, "\nwait: %d; max: %d; size: %d\n", wait, max, n)
		if err != nil {
			return n, err
		}
		time.Sleep(time.Duration(wait) * time.Millisecond)
		return n, err
	})
}
