- https://github.com/gookit/config
- https://github.com/golang-migrate/migrate
- https://pkg.go.dev/log/slog

# nsqsink washtub

## Getting Started

TBD

## How to Use 

TBD

### Libraries & Tools Used

TBD

### Folder Structure

Here is the folder structure we have been using in this project

```
cmd/
|- main.go
docker/
|- docker-compose.yml
pkg/
server/
|- domain/
|- handler/
|- repository/
|- usecase/
|- http.go
static/
|- build/
|- css/
|- fonts/
|- html/
|- js/
```

Now, lets dive into each folder.

```
1- cmd - The initiation of go binary run.
2- docker - Contains the docker setups.
3- pkg - Contains the package/utilities/common functions.
4- server — Contains all the backend server for washtub
5- static — Contains all the frontend and assets for washtub
```

### To Do

- API Contract
- Server Init
- Frontend Init