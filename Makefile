NS = kjbreil
REPO = goscript
VERSION ?= $(shell date +'%Y.%m.%d')

.PHONY: image no-cache push

default: image push

image:
	go mod vendor
	docker build -t $(NS)/$(REPO) -t $(NS)/$(REPO):$(VERSION) .

no-cache:
	go mod vendor
	docker build --no-cache -t $(NS)/$(REPO) -t $(NS)/$(REPO):$(VERSION) .

run:
	docker run --rm --name goscript $(NS)/$(REPO):$(VERSION)

start:
	docker run -d --name goscript $(NS)/$(REPO):$(VERSION)

stop:
	docker stop goscript

rm:
	docker rm goscript
