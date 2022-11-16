# Rest api for Video Platform

## Postgres

#### create docker postgres
- docker run --name=postgres -e POSTGRES_PASSWORD='qwerty' -p 5432:5432 -d postgres

#### create init database
- migrate create -ext sql -dir ./migrations -seq init

#### create migration. [golang-migrate/migrate](https://github.com/golang-migrate/migrate)
- migrate -path ./migrations -database "postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable" up

## Redis

#### create docker redis
- docker run --name=redis -p 6379:6379 -d redis