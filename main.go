package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/koki-develop/go-fzf"
)

const version = "0.0.1"

const appDir = ".inventory"

var revision = "HEAD"

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	os.Exit(1)
}

func dirWalk(base string) ([]string, error) {
	var files []string
	err := filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(base, path)
		if err != nil {
			return err
		}
		files = append(files, rel)
		return nil
	})

	if err != nil {
		return nil, err
	}
	return files, nil
}

func copy(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	var showVersion bool

	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.Parse()

	if showVersion {
		fmt.Printf("version: %s, revision: %s\n", version, revision)
		return
	}

	if flag.NArg() > 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s filename\n", os.Args[0])
		os.Exit(1)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fatal(err)
	}
	files, err := dirWalk(filepath.Join(homeDir, appDir))
	if err != nil {
		fatal(err)
	}

	f, err := fzf.New()
	if err != nil {
		fatal(err)
	}

	idxs, err := f.Find(files, func(i int) string { return files[i] })
	if err != nil {
		fatal(err)
	}

	src := filepath.Join(homeDir, appDir, files[idxs[0]])

	wd, err := os.Getwd()
	if err != nil {
		fatal(err)
	}

	dst := filepath.Join(wd, filepath.Base(src))

	err = copy(src, dst)
	if err != nil {
		fatal(err)
	}

}
