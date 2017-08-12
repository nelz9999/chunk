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

var cfgFile string

var debug bool
var maxSize int
var lowSize int
var maxWait int
var lowWait int
var input string

// var output string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "chunk",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		in := buildReader()
		defer in.Close()

		out := buildWriter()
		defer out.Close()

		log := buildLog()

		sizer, buf := buildSizer()
		waiter := buildWaiter()
		src := stream.New(in, sizer, waiter, log)

		io.CopyBuffer(out, src, buf)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.chunk.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "send debugging output to stderr")
	RootCmd.Flags().IntVarP(&maxSize, "max-size", "s", 16, "set the maximum chunk size to send")
	RootCmd.Flags().IntVarP(&lowSize, "low-size", "l", 0, "set to a non-zero value less than the max-size to send random variable sized chunks")
	RootCmd.Flags().IntVarP(&maxWait, "max-wait", "w", 100, "set the period, in milliseconds, to wait between chunk delivery")
	RootCmd.Flags().IntVarP(&lowWait, "min-wait", "m", 0, "set to a non-zero value less than the max-wait to wait random variable periods between chunks")
	RootCmd.Flags().StringVarP(&input, "input", "i", "", "specify source file, otherwise defaults to stdin")
	// RootCmd.Flags().StringVarP(&output, "ouput", "o", "", "specify destination file, otherwise defaults to stdout")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// if cfgFile != "" {
	// 	// Use config file from the flag.
	// 	viper.SetConfigFile(cfgFile)
	// } else {
	// 	// Find home directory.
	// 	home, err := homedir.Dir()
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		os.Exit(1)
	// 	}
	//
	// 	// Search config in home directory with name ".chunk" (without extension).
	// 	viper.AddConfigPath(home)
	// 	viper.SetConfigName(".chunk")
	// }
	//
	// viper.AutomaticEnv() // read in environment variables that match
	//
	// // If a config file is found, read it in.
	// if err := viper.ReadInConfig(); err == nil {
	// 	fmt.Println("Using config file:", viper.ConfigFileUsed())
	// }
}

func buildReader() io.ReadCloser {
	result := os.Stdin

	if input != "" {
		var err error
		result, err = os.Open(input)
		if err != nil {
			panic(err)
		}
	}
	return result
}

func buildWriter() io.WriteCloser {
	result := os.Stdout

	// if output != "" {
	// 	var err error
	// 	result, err = os.OpenFile(
	// 		output,
	// 		os.O_WRONLY|os.O_APPEND|os.O_CREATE,
	// 		0600)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	return result
}

func buildLog() io.Writer {
	result := ioutil.Discard
	if debug {
		result = os.Stderr
	}
	return result
}

func buildSizer() (stream.Inter, []byte) {
	buf := make([]byte, maxSize)
	sizer := stream.InterFunc(func() int {
		return maxSize
	})

	if lowSize > maxSize {
		panic(fmt.Sprintf("lowSize > maxSize: %d > %d", lowSize, maxSize))
	}

	if lowSize > 0 && lowSize < maxSize {
		breadth := maxSize - lowSize + 1
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		sizer = stream.InterFunc(func() int {
			return lowSize + r.Intn(breadth)
		})
	}

	return sizer, buf
}

func buildWaiter() stream.Inter {
	waiter := stream.InterFunc(func() int {
		return maxWait
	})

	if lowWait > maxWait {
		panic(fmt.Sprintf("lowWait > maxWait: %d > %d", lowWait, maxWait))
	}

	if lowWait > 0 && lowWait < maxWait {
		breadth := maxWait - lowWait + 1
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		waiter = stream.InterFunc(func() int {
			return lowWait + r.Intn(breadth)
		})
	}

	return waiter
}
