all: fmt vet install test

fmt:
	go fmt ./...

vet:
	go vet ./...

install:
	go install ./...

test:
	go test ./...

	~/go/bin/server -help

dumpSchema:
	mysqldump \
		--protocol=tcp \
		--user=${MYSQL_USER} \
		--compact \
		--comments \
		--no-data \
		--databases Water \
	| sed -e 's/^CREATE TABLE /CREATE TABLE IF NOT EXISTS /' \
	> db/mysqldb/schema.sql

eraseSchema:
	echo 'drop database Water' | mysql Water

