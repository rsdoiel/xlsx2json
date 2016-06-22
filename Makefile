#
# Simple Makefile for conviently testing, building and deploying experiment.
#
build:
	go build -o bin/xlsx2json cmds/xlsx2json/xlsx2json.go

test:
	#go test
	gocyclo -over 14 .

clean:
	if [ -d bin ]; then rm -fR bin; fi
	if [ -d dist ]; then rm -fR dist; fi

install:
	GOBIN=$HOME/bin go install cmds/xlsx2json/xlsx2json.go

release:
	./mk-release.sh


