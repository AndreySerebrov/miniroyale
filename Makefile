test:
	go clean -testcache
	go test -v ./...

docker-run:
	docker build . -t pow
	docker run pow