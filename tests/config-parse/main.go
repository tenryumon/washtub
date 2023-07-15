// This is a hack of config template parsing.
//
// WARN: This test is only a best effort test.

package main

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/google/go-cmp/cmp"
	"github.com/nyelonong/boilerplate-go/core/config"
	testpkg "github.com/nyelonong/boilerplate-go/pkg/testing"
)

func main() {
	topLevel, err := testpkg.RepoTopLevel()
	if err != nil {
		panic(err)
	}

	// Set all environment variables.
	os.Setenv("MYSQL_PASSWORD", "testing")
	os.Setenv("MAILER_PASSWORD", "testing")
	os.Setenv("WHATSAPP_TOKEN", "testing")
	os.Setenv("QISCUS_NAMESPACE", "testing")
	os.Setenv("QISCUS_APP_ID", "testing")
	os.Setenv("QISCUS_SECRET_KEY", "testing")
	os.Setenv("QISCUS_TOGGLE", "testing")
	os.Setenv("XENDIT_SECRET_KEY", "testing")

	configName := "backend-PRODUCTION.ini"
	configToTest := path.Join(topLevel, "devops/configuration", configName)

	out, err := config.ParseFile(configToTest)
	if err != nil {
		panic(err)
	}

	goldenFile := "backend-production.golden"
	gf, err := os.Open(goldenFile)
	if err != nil {
		panic(err)
	}
	outGolden, err := io.ReadAll(gf)
	if err != nil {
		panic(err)
	}

	if diff := cmp.Diff(string(outGolden), string(out)); diff != "" {
		fmt.Println(diff)
		panic("diff output")
	}
	fmt.Println("OK")
}
