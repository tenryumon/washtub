package config

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadFile(t *testing.T) {
	filepath := "testdata/tdata1.ini"

	t.Setenv("MYSQL_USERNAME", "user")
	t.Setenv("MYSQL_PASSWORD", "mysql")

	conf := struct {
		Database struct {
			Driver     string
			Connection string
		}
	}{}
	expect := struct {
		Database struct {
			Driver     string
			Connection string
		}
	}{
		Database: struct {
			Driver     string
			Connection string
		}{
			Driver:     "mysql",
			Connection: "user:mysql@tcp(mysql:3306)/dbname?parseTime=true",
		},
	}

	if err := ReadFile(&conf, filepath); err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expect, conf); diff != "" {
		t.Fatalf("(-want/+got)\n%s", diff)
	}
}
