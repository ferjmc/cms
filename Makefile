.PHONY: build clean deploy

build:
	chmod u+x gobuild.sh
	./gobuild.sh

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose
