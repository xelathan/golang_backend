#!/bin/bash

migrate -path ~/golang/cmd/migrate/migrations -database "mysql://root:mypassword@tcp(127.0.0.1:3306)/golang_db" force $NEW_VERSION