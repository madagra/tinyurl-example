BINARY_NAME=app
APP_NAME=tinyurl

compile:
	go build -C ${APP_NAME} -o `pwd`/${BINARY_NAME}

run: compile
	`pwd`/${BINARY_NAME} --local

test: compile
	docker compose up -d
	go test -v `pwd`/${APP_NAME}
	docker compose down --volumes

clean:
	go clean
	rm `pwd`/${BINARY_NAME}

build_image:
	docker build -t ${APP_NAME}:dev .

push_image: build_image
	echo ${DOCKER_USR} ${DOCKER_PWD}
	docker tag ${APP_NAME}:dev madagra/${APP_NAME}:dev
	docker login -u ${DOCKER_HUB_USERNAME} -p ${DOCKER_HUB_PASSWORD}
	docker push madagra/${APP_NAME}:dev

run_image: build_image
	docker run -p 3000:3000 ${APP_NAME}:dev 
