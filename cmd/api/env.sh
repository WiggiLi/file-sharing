#!bin/bash
export HTTP_ADDR=localhost:8888

export FILE_PATH=../../files

export PG_URL=postgres://postgres:mypass@localhost/postgres?sslmode=disable
export PG_MIGRATIONS_PATH=file://../../store/pg/migrations

export LOG_LEVEL=debug