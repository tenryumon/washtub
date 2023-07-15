package testdb

import (
	"path"
	"testing"

	testingpkg "github.com/nyelonong/boilerplate-go/pkg/testing"
)

func TestMigrateUp(t *testing.T) {
	tdb, err := New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		err := tdb.Drop()
		if err != nil {
			t.Log(err)
		}
	})

	migrationFolder, err := testingpkg.RepoTopLevel()
	if err != nil {
		t.Fatal(err)
	}
	migrationFolder = path.Join(migrationFolder, "sqlfiles")

	if err := tdb.MigrateUP(migrationFolder); err != nil {
		t.Fatal(err)
	}
}
