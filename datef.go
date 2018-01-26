package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"
)

func main() {
	var (
		iFormatStr = flag.String("i", "unix", "Input `format`")
		oFormatStr = flag.String("o", "RFC3339", "Output `format`")
	)
	flag.Usage = usage
	flag.Parse()

	iformat := newFormat(*iFormatStr)
	oformat := newFormat(*oFormatStr)

	if flag.NArg() == 0 {
		fmt.Println(oformat.format(time.Now()))
		return
	}
	if flag.NArg() == 1 && flag.Arg(0) == "-" {
		ok := true
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			s := scanner.Text()
			if s == "" {
				continue
			}
			if !printTimestamp(s, iformat, oformat) {
				ok = false
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			ok = false
		}
		if !ok {
			os.Exit(1)
		}
		return
	}
	ok := true
	for _, s := range flag.Args() {
		if !printTimestamp(s, iformat, oformat) {
			ok = false
		}
	}
	if !ok {
		os.Exit(1)
	}
}

func printTimestamp(s string, iformat, oformat format) (ok bool) {
	t, err := iformat.parse(s)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return false
	}
	t = t.UTC()
	fmt.Println(oformat.format(t))
	return true
}

type formatType int

const (
	typeUnix formatType = iota
	typeUnixMS
	typeNamed
	typeGoFormat
)

type format struct {
	typ formatType
	s   string
}

type namedFormat struct {
	f    string
	desc string
}

var namedFormats = map[string]namedFormat{
	"RFC3339": {time.RFC3339, "RFC3339 timestamp"},
}

func newFormat(s string) format {
	switch s {
	case "unix":
		return format{typeUnix, ""}
	case "unixms":
		return format{typeUnixMS, ""}
	}
	if _, ok := namedFormats[s]; ok {
		return format{typeNamed, s}
	}
	return format{typeGoFormat, s}
}

func (f format) String() string {
	switch f.typ {
	case typeUnix:
		return "unix"
	case typeUnixMS:
		return "unixms"
	case typeNamed, typeGoFormat:
		return f.s
	}
	panic("unreached")
}

func (f format) parse(s string) (time.Time, error) {
	var (
		t   time.Time
		err error
		ok  = true
	)
	switch f.typ {
	case typeUnix, typeUnixMS:
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			ok = false
			break
		}
		if f.typ == typeUnix {
			t = time.Unix(n, 0)
		} else {
			t = time.Unix(0, n*1000*1000)
		}
	case typeNamed:
		t, err = time.Parse(namedFormats[f.s].f, s)
		if err != nil {
			ok = false
		}
	case typeGoFormat:
		t, err = time.Parse(f.s, s)
		if err != nil {
			ok = false
		}
	}
	if ok {
		return t, nil
	}
	return t, fmt.Errorf("%q is an invalid value for format %s", s, f)
}

func (f format) format(t time.Time) string {
	switch f.typ {
	case typeUnix:
		return strconv.FormatInt(t.Unix(), 10)
	case typeUnixMS:
		return strconv.FormatInt(t.UnixNano()/1000/1000, 10)
	case typeNamed:
		return t.Format(namedFormats[f.s].f)
	case typeGoFormat:
		return t.Format(f.s)
	}
	panic("unreached")
}

func usage() {
	fmt.Fprintf(os.Stderr, `usage:

  %s [flags]                            or
  %[1]s [flags] -                          or
  %[1]s [flags] timestamp1 timestamp2 ...

where the flags are:

`, os.Args[0])
	flag.PrintDefaults()
	fmt.Fprint(os.Stderr, `
and the formats are:

  unix     seconds since unix epoch
  unixms   milliseconds since unix epoch
`)

	var names []string
	for name := range namedFormats {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		nf := namedFormats[name]
		fmt.Fprintf(os.Stderr, "  %-8s %s\n", name, nf.desc)
	}

	fmt.Fprint(os.Stderr, `
or any custom timestamp format as permitted by https://golang.org/pkg/time.

One or more timestamps may be provided; the output will be printed line-by-line.
If no timestamps are provided, the current time is used. If - is given as the
only argument, then input timestamps are read line-by-line from standard input.
`)
}
