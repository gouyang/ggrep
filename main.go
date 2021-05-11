package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

var colors = map[string]string{
	"reset":  "\033[0m",
	"red":    "\033[31m",
	"green":  "\033[32m",
	"yellow": "\033[33m",
	"blue":   "\033[34m",
	"purple": "\033[35m",
	"cyan":   "\033[36m",
	"white":  "\033[37m",
}

func search(reader *bufio.Reader, s string, invert bool, color string) {
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			switch err {
			default:
				fmt.Errorf("unable to read: %s", err)
			case io.EOF:
				return
			}
		}

		content := string(line)
		re := regexp.MustCompile(s)
		// invert-matchings
		if invert && !re.MatchString(content) {
			fmt.Println(content)
		}
		// matchings
		if re.MatchString(content) && !invert {
			for {
				index := re.FindIndex(line)
				if index != nil {
					fmt.Print(string(line[:index[0]]))
					fmt.Print(string(colors[color]), string(line[index[0]:index[1]]), colors["reset"])
					line = line[index[1]:]
				} else {
					fmt.Println(string(line))
					break
				}
			}
		}
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: ggrep searchstring inputfile\n")
	flag.PrintDefaults()
	os.Exit(2)
}

var invert bool
var color string

func init() {
	flag.BoolVar(&invert, "v", false, "select non-matching lines")
	flag.StringVar(&color, "c", "red", "color for matching string, support colors: red, green, yellow, purple, cyan, white, blue")
	// TODO: -r, --recursive
	// TODO: -i, --ignore-case
	// TODO: -n, --line-number
}

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()

	// handle piper. If the stdin is pipe, execute and exit.
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	if fi.Mode()&os.ModeNamedPipe != 0 {
		reader := bufio.NewReader(os.Stdin)
		str := strings.Join(args, " ")
		search(reader, str, invert, color)
		os.Exit(0)
	}

	if len(args) < 2 {
		fmt.Println("Input file is missing.")
		os.Exit(1)
	}

	for _, file := range args[1:] {
		f, err := os.Open(file)
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()
		buf := bufio.NewReader(f)
		search(buf, args[0], invert, color)
	}
}
