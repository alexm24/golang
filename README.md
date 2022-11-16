# Rest api for Video Platform

## Swagger. OpenAPI 3.0
- https://github.com/alexm24/golang/blob/main/internal/handler/api/api.swagger.yaml

## Postgres

#### create docker postgres
- docker run --name=postgres -e POSTGRES_PASSWORD='qwerty' -p 5432:5432 -d postgres

#### create init database
- migrate create -ext sql -dir ./migrations -seq init

#### create migration with [golang-migrate](https://github.com/golang-migrate/migrate)
- migrate -path ./migrations -database "postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable" up

## Redis

#### create docker redis
- docker run --name=redis -p 6379:6379 -d redis

## Centrifugo
#### [Centrifugo is an open-source scalable real-time messaging server.](https://github.com/centrifugal/centrifugo)
- docker run --ulimit nofile=65536:65536 -v /host/dir/with/config/file:/centrifugo -p 8000:8000 centrifugo/centrifugo centrifugo -c config.json