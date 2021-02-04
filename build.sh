#!/bin/sh
git pull --ff
go build -o ./bin/ ./cmd/bot
