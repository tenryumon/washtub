export TESTING_DB_DSN=user:asdasd@tcp(localhost:3306)/dbname?parseTime=true
export TESTING_DB_DRIVER=mysql
export TESTING_DB_NAME=dbname

test:
	@go test -v -race `go list ./... | grep -e usecases/billing -e mysql` -p=1

.PHONY: test
