package testdb

import (
	"context"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/nyelonong/boilerplate-go/core/database"
)

// TESTING_DB_DRIVER is the database driver that we use for testing. We can use multiple database drivers
// based on the DSN.
var dbDriver = os.Getenv("TESTING_DB_DRIVER")

// TESTING_DB_DSN is the data source name that we use for testing. The DSN determines the database connection.
var dbDSN = os.Getenv("TESTING_DB_DSN")

// TESTING_DB_NAME is the database name that we use for testing. The database name is needed because we need to
// drop the database when the tests is over.
var dbName = os.Getenv("TESTING_DB_NAME")

type TestDB struct {
	DB *database.DB

	driver string
	dsn    string
	dbname string
}

// New craetes a new test db object, re-create the database and also connect the test db itself to a database.
func New() (*TestDB, error) {
	if err := recreateDB(dbDriver, dbDSN); err != nil {
		return nil, err
	}

	db, err := database.Connect(dbDriver, dbDSN)
	if err != nil {
		return nil, err
	}
	return &TestDB{
		DB:     db,
		driver: dbDriver,
		dsn:    dbDSN,
		dbname: dbName,
	}, nil
}

func recreateDB(driver, dsn string) error {
	conf, err := parseMYSQLDSN(dsn)
	if err != nil {
		return err
	}
	createDBDSN := fmt.Sprintf("%s:%s@tcp(%s)/", conf.User, conf.Passwd, conf.Addr)

	db, err := database.Connect(driver, createDBDSN)
	if err != nil {
		return err
	}

	dropDBQuery := fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName)
	_, err = db.Exec(context.Background(), dropDBQuery, nil)
	if err != nil {
		return err
	}

	createDBQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName)
	_, err = db.Exec(context.Background(), createDBQuery, nil)
	return err
}

func (t *TestDB) MigrateUP(sqlFolder string) error {
	db, err := database.Connect(dbDriver, dbDSN)
	if err != nil {
		return err
	}
	defer db.Close()

	sqlFiles, err := getAllFile(sqlFolder)
	if err != nil {
		return err
	}

	// Set the migration deadline to 30s, we can always adjust this later on.
	migrationTimeout, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	for _, f := range sqlFiles {
		err = readAndRun(migrationTimeout, db, f)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TestDB) Drop() error {
	// Enforce the connection close when dropping the database.
	t.DB.Close()

	// We will ignore the driver for now as we only use MySQL.
	conf, err := parseMYSQLDSN(t.dsn)
	if err != nil {
		return err
	}
	dropDSN := fmt.Sprintf("%s:%s@tcp(%s)/", conf.User, conf.Passwd, conf.Addr)

	// Reconnect using the provided dsn information.
	db, err := database.Connect(t.driver, dropDSN)
	if err != nil {
		return fmt.Errorf("drop_db: failed to connect to database. error: %v", err)
	}
	defer db.Close()

	// if err := killAllMySQLProcess(db, conf.User); err != nil {
	// 	return err
	// }

	query := fmt.Sprintf("DROP DATABASE IF EXISTS %s", t.dbname)
	_, err = db.Exec(context.Background(), query, nil)
	if err != nil {
		err = fmt.Errorf("drop_db: failed to drop database %s with error: %v", t.dbname, err)
	}
	return err
}

// killAllMySQLProcesses kill all remaining processes so we can drop the MySQL database.
func killAllMySQLProcess(db *database.DB, user string) error {
	query := fmt.Sprintf("select concat('KILL ',id,';') from information_schema.processlist where user='%s';", user)

	killCmds := []string{}
	if err := db.Select(context.Background(), &killCmds, query, nil); err != nil {
		return fmt.Errorf("kill_mysql_processes: failed to retrieve all pricesses on user %s with error: %v", user, err)
	}

	for _, killCmd := range killCmds {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		_, err := db.Exec(ctxTimeout, killCmd, nil)
		if err != nil {
			return fmt.Errorf("kill_mysql_process: failed to kill process with query: %s and error: %v", killCmd, err)
		}
	}
	return nil
}

func getAllFile(folderLoc string) ([]string, error) {
	folder, err := os.Open(folderLoc)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s because %s", folderLoc, err)
	}

	files, err := folder.Readdir(0)
	if err != nil {
		return nil, fmt.Errorf("failed to get file list in %s because %s", folderLoc, err)
	}

	list := []string{}
	for _, v := range files {
		if v.IsDir() {
			continue
		}
		if !strings.HasSuffix(v.Name(), ".sql") {
			continue
		}
		list = append(list, path.Join(folderLoc, v.Name()))
	}
	sort.Strings(list)
	return list, nil
}

func readAndRun(ctx context.Context, db *database.DB, fileLoc string) error {
	content, err := os.ReadFile(fileLoc)
	if err != nil {
		return fmt.Errorf("failed to read file because %s", err)
	}

	queries := string(content)
	if queries == "" {
		return fmt.Errorf("failed to read file because file is empty")
	}

	for _, query := range strings.Split(queries, ";") {
		_, err = db.Exec(ctx, query, nil)
		if err != nil {
			// skip empty query
			if err.Error() == "Error 1065: Query was empty" {
				continue
			}
			return fmt.Errorf("failed to execute sql:\n%s\n\non file %s because %s", query, fileLoc, err)
		}
	}

	return nil
}
