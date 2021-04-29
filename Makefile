deploy-env:
ifndef NAME
	$(error NAME is not defined)
endif
ifndef MOUNT
	$(error MOUNT is not defined)
endif

deploy: deploy-env
	docker build -t cxevolves/$(NAME) . 
	docker run -it --rm --name cxevolves-run -v $(MOUNT):/Results cxevolves/$(NAME):latest