package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/nyelonong/boilerplate-go/core/config"
	"github.com/nyelonong/boilerplate-go/core/database"
	"github.com/nyelonong/boilerplate-go/core/environment"
	"github.com/nyelonong/boilerplate-go/core/log"
)

type Configuration struct {
	Database DatabaseConfig
	Log      log.Config
}
type DatabaseConfig struct {
	Driver     string
	Connection string
}

const separator = ";"

func readAndRun(ctx context.Context, db *database.DB, fileLoc string) error {
	content, err := os.ReadFile(fileLoc)
	if err != nil {
		return fmt.Errorf("Failed to read file because %s", err)
	}

	queries := string(content)
	if queries == "" {
		return fmt.Errorf("Failed to read file because file is empty")
	}

	for _, query := range strings.Split(queries, separator) {
		println(query)
		_, err = db.Exec(ctx, query, nil)
		if err != nil {
			// skip empty query
			if err.Error() == "Error 1065: Query was empty" {
				continue
			}
			return fmt.Errorf("Failed to execute sql file because %s", err)
		}
		println("SUCCESS")
	}

	return nil
}

func getAllFile(folderLoc string, prefix string) ([]string, error) {
	folder, err := os.Open(folderLoc)
	if err != nil {
		return nil, fmt.Errorf("Failed to open %s because %s", folderLoc, err)
	}

	files, err := folder.Readdir(0)
	if err != nil {
		return nil, fmt.Errorf("Failed to get file list in %s because %s", folderLoc, err)
	}

	list := []string{}
	for _, v := range files {
		if v.IsDir() {
			continue
		}
		if !strings.HasSuffix(v.Name(), ".sql") {
			continue
		}
		if prefix != "" && !strings.HasPrefix(v.Name(), prefix) {
			continue
		}
		list = append(list, folderLoc+v.Name())
	}
	sort.Strings(list)
	return list, nil
}

func main() {
	// Get Flag for service
	var configFile string
	var sqlFile string
	var sqlFolder string
	var sqlPrefix string
	var runAll bool
	configPath := "devops/configuration/backend-development.ini"
	if environment.IsStaging() {
		configPath = "devops/configuration/backend-staging.ini"
	}
	flag.StringVar(&configFile, "config", configPath, "Configuration File Location")
	flag.StringVar(&sqlFile, "sqlfile", "sqlfiles/001-init.sql", "SQL File to run location")
	flag.StringVar(&sqlFolder, "sqlfolder", "sqlfiles/", "Sql Files Folder")
	flag.StringVar(&sqlPrefix, "sqlprefix", "", "Filter Sql Files by prefix")
	flag.BoolVar(&runAll, "all", false, "Run all sqlfiles")
	flag.Parse()

	ctx := context.Background()

	// Read and parse config
	conf := Configuration{}
	err := config.ReadFile(&conf, configFile)
	if err != nil {
		log.Fatalf("Failed to read configuration because %s", err)
	}

	// Initialize Log
	log.Init(conf.Log)

	// Connect to Database Master
	masterDB, err := database.Connect(conf.Database.Driver, conf.Database.Connection)
	if err != nil {
		log.Fatalf("Failed to connect to database because %s", err)
	}

	fileList := []string{}
	if runAll {
		fileList, err = getAllFile(sqlFolder, sqlPrefix)
		if err != nil {
			log.Fatalf("Failed because %s", err)
		}
	} else {
		fileList = append(fileList, sqlFile)
	}

	for _, fileLoc := range fileList {
		log.Infof("-- [START] %s", fileLoc)
		err = readAndRun(ctx, masterDB, fileLoc)
		if err != nil {
			log.Fatalf("Failed because %s", err)
		}
		log.Infof("-- [SUCCESS] %s", fileLoc)
	}

	log.Infof("Success run all sql files")
}
