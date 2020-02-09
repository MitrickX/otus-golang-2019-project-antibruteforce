.PHONY: run stop restart test build

run:
	docker-compose -f ./build/package/docker-compose.yml up -d --build

stop:
	docker-compose -f ./build/package/docker-compose.yml down

restart: stop run

build:
	go build -o antibruteforce .

test:
	set -e ;\
	tests_status_code=0 ;\
	docker-compose -f ./build/package/docker-compose-tests.yml down ;\
	docker-compose -f ./build/package/docker-compose-tests.yml up -d --build ;\
	docker-compose -f ./build/package/docker-compose-tests.yml run tests ./ip-tests && \
		docker-compose -f ./build/package/docker-compose-tests.yml run tests ./grpc-tests --config ./configs/config.yml --features-path ./features/ || \
			tests_status_code=$$? ;\
	exit $$tests_status_code ;\
