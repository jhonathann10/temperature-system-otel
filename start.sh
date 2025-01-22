#!/bin/sh
go run ./cmd/api/main.go &
go run ./cmd/apib/main.go &
wait