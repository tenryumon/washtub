package main

import (
	"context"
	"flag"

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

type User struct {
	ID    int64  `db:"id"`
	Name  string `db:"name"`
	Phone string `db:"phone"`
	Email string `db:"email"`
}

func main() {
	// Get Flag for service
	var configFile string
	configPath := "devops/configuration/backend-development.ini"
	if environment.IsStaging() {
		configPath = "devops/configuration/backend-staging.ini"
	}
	flag.StringVar(&configFile, "config", configPath, "Configuration File Location")
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

	query := `CREATE TABLE dummies (
	id             	VARCHAR(36)		NOT NULL,
	bigint   		BIGINT       	NOT NULL,
	text           	VARCHAR(255) 	NOT NULL,
	created_by     	BIGINT       	NOT NULL,
	created_time   	TIMESTAMP    	NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_by     	BIGINT       	NOT NULL,
	updated_time  	TIMESTAMP    	NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY(id)
);`

	_, err = masterDB.Exec(ctx, query, nil)
	if err != nil {
		log.Fatalf("Failed to query because %s", err)
	}

	log.Infof("Success run all sql files")
}
