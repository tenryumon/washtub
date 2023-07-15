# Go Backend Boilerplate

## Deps
---

Backend requires Golang and Docker to run. Please install from below url for your local machine.
- Golang
- Docker

Start dependencies using docker
```sh
# Go to your service root
# Setup Dependecies and automatically insert database in sqlfiles.
# (Only the first initialize, to reset database, you need to remove volume first)
make docker-start

# If you want to stop the docker process, do below:
make docker-stop

# If you want to remove docker image, this won't reset your volume data
make docker-remove
```

Start the service by run functions from Makefile
```sh
# Go to your service root
# Start the service
make run
```

When there is an additional changes in database after your first docker-start, you can run the additional sqlfile
```sh
# Go to your service root
# Check which sqlfiles you need to run
ls sqlfiles/

# Get the 3 first digit of the sqlfile (example is 011)
make run-db prefix=011
````

Postman collection is inside `postman` folder.
It will be updated periodically, import the newest postman json to your postman application

## Folder Structure
---

```sh
root/
|--- cmd
|--- core
|--- devops
|--- docker
|--- internal
     |--- handlers
     |--- usecases
     |--- repositories
     |--- models
     |--- entities
     |--- interfaces
```