package main

import (
	"flag"
	"os"
	"path"

	"github.com/nyelonong/boilerplate-go/core/config"
	"github.com/nyelonong/boilerplate-go/pkg/testing"
)

func main() {
	var (
		p  string
		tp string
	)
	flag.StringVar(&p, "config-path", "", "parse configuration file")
	flag.StringVar(&tp, "target-path", "", "create new configuration file")
	flag.Parse()

	// This is only a HACK, please don't use testing package directly on production.
	topLevel, err := testing.RepoTopLevel()
	if err != nil {
		panic(err)
	}
	// Join the relative path of configuration to the repository top path.
	p = path.Join(topLevel, p)

	out, err := config.ParseFile(p)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(tp)
	if err != nil {
		panic(err)
	}
	_, err = f.Write(out)
	if err != nil {
		panic(err)
	}
	f.Close()
}
