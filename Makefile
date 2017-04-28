NAME = fabriziopandini/test-service
VERSION = 0.1

all: package

package:
	@docker run --rm -v $(PWD):/src -v /var/run/docker.sock:/var/run/docker.sock fabriziopandini/golang-builder $(NAME):$(VERSION) 

test_package: 
	@docker run -p 8080:8080 --rm $(NAME):$(VERSION)

tag: 
	@docker tag $(NAME):$(VERSION) $(NAME):latest 
    
push: 
	@docker push $(NAME)
	
rmi: 
	@docker rmi $(NAME):$(VERSION) $(NAME):latest 