#!/bin/bash
go list -f '{{range .Imports}}{{.}} {{end}}' ./escher.go | xargs go get -v
go get golang.org/x/tools/cmd/cover
