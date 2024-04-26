// Copyright 2021 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Mdfmt reformats Markdown data.
//
// Usage:
//
//	mdfmt [-w] [file...]
//
// Mdfmt reads the named files, or else standard input, as Markdown documents
// and then reprints the same Markdown documents to standard output.  Includes
// support for the extended tables syntax.
//
// The -w flag specifies to rewrite the files in place.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/debug"

	markdown "github.com/zacharysyoung/rsc_markdown"
)

var (
	wflag = flag.Bool("w", false, "write reformatted Markdown to files ")
	vflag = flag.Bool("v", false, "version info")
	exit  = 0
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: mdfmt [-v] [-w] [file...]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func version() {
	s := "zachs_markdown:mdfmt"
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		panic(fmt.Errorf("could not get build info"))
	}
	for _, x := range bi.Settings {
		if x.Key == "vcs.revision" {
			s += ":" + x.Value[:7] // short hash
			break
		}
	}
	s += ":" + bi.GoVersion
	fmt.Fprintln(os.Stderr, s)
	os.Exit(2)
}

func main() {
	log.SetPrefix("mdfmt: ")
	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	if *vflag {
		version()
	}

	if flag.NArg() == 0 {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		convert(data, "")
	} else {
		for _, file := range flag.Args() {
			data, err := os.ReadFile(file)
			if err != nil {
				log.Print(err)
				exit = 1
				continue
			}
			convert(data, file)
		}
	}
	os.Exit(exit)
}

func convert(data []byte, file string) {
	p := markdown.Parser{Table: true}
	doc := p.Parse(string(data))
	out := []byte(markdown.ToMarkdown(doc))
	if *wflag && file != "" {
		if err := os.WriteFile(file, out, 0666); err != nil {
			log.Print(err)
			exit = 1
			return
		}
	} else {
		os.Stdout.Write(out)
	}
}
