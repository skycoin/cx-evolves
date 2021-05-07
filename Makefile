build-env:
ifndef NAME
	$(error NAME is not defined)
endif

deploy-env: build-env
ifndef MOUNT
	$(error MOUNT is not defined)
endif

build: build-env
	go mod tidy
	go mod download
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build .
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build ./cxexecutes/worker/cmd/server.go
	docker build -t cxevolves/$(NAME) . 
	rm cx-evolves
	rm server

deploy:build deploy-env
	docker run -it --rm --name cxevolves-run -v $(MOUNT):/Results cxevolves/$(NAME):latest

push-docker: build-env build
	docker tag cxevolves/$(NAME):latest kenje4090/kenjecxevolves 
	docker push kenje4090/kenjecxevolves:latest