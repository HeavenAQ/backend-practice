# Docker Settings
SYS_PORT=5432
DOCKER_PORT=5432
DOCKER_CONTAINER=backend-postgres
DOCKER_IMG=postgres
DOCKER_IMG_VER=12-alpine

# Database Settings
USER=root
PASSWD=Computationalinguist
DATABASE=simple_bank
DATABASE_PATH=db/migration

postgres:
	docker run --name ${DOCKER_CONTAINER} -p ${SYS_PORT}:${DOCKER_PORT} -e POSTGRES_USER=${USER} -e POSTGRES_PASSWORD=${PASSWD} -d ${DOCKER_IMG}:${DOCKER_IMG_VER}

startdb:
	docker container start ${DOCKER_CONTAINER}

createdb:
	docker exec -it ${DOCKER_CONTAINER} createdb --username=${USER} --owner=${USER} ${DATABASE}

migrateup:
	migrate -path ${DATABASE_PATH} -database "postgresql://${USER}:${PASSWD}@localhost:${SYS_PORT}/${DATABASE}?sslmode=disable" -verbose up

migratedown:
	migrate -path ${DATABASE_PATH} -database "postgresql://${USER}:${PASSWD}@localhost:${SYS_PORT}/${DATABASE}?sslmode=disable" -verbose down

test:
	go test -v -cover ./...

dropdb:
	docker exec -it ${DOCKER_CONTAINER} dropdb --username=${USER} ${DATABASE}

sqlc:
	sqlc generate

server:
	go run main.go

.PHONY: createdb dropdb postgres migrateup migratedown sqlc server
