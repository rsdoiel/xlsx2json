#
# Simple Makefile for conviently testing, building and deploying experiment.
#
build:
	go build -o bin/xlsx2json cmds/xlsx2json/xlsx2json.go

test:
	#go test
	gocyclo -over 14 .

clean:
	if [ -f bin/xlsx2json ]; then rm bin/xlsx2json; fi

install:
	go install cmds/xlsx2json/xlsx2json.go
