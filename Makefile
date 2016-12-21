#
# Simple Makefile for conviently testing, building and deploying experiment.
#
PROJECT = xlsx2json

VERSION = $(shell grep -m 1 'Version =' $(PROJECT).go | cut -d\"  -f 2)

BRANCH = $(shell git branch | grep '* ' | cut -d\  -f 2)

build:
	env CGO_ENABLED=0 go build -o bin/xlsx2json cmds/xlsx2json/xlsx2json.go

test:
	#go test
	gocyclo -over 14 .

save:
	git commit -am "Quick Save"
	git push origin $(BRANCH)

clean:
	if [ -d bin ]; then /bin/rm -fR bin; fi
	if [ -d dist ]; then /bin/rm -fR dist; fi
	if [ -f $(PROJECT)-$(VERSION)-release.zip ]; then /bin/rm $(PROJECT)-$(VERSION)-release.zip; fi

install:
	env CGO_ENABLED=0 GOBIN=$(HOME)/bin go install cmds/xlsx2json/xlsx2json.go

release:
	./mk-release.sh


