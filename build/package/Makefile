.PHONY: up down restart test

up:
	docker-compose -f docker-compose-tests.yml up -d --build

dowm:
	docker-compose down

restart: down up


test:
	set -e ;\
	tests_status_code=0 ;\
	docker-compose down ;\
	docker-compose -f docker-compose-tests.yml up -d --build ;\
	docker-compose -f docker-compose-tests.yml run tests ./tests --config ./configs/config.yaml --features-path ./features/ || tests_status_code=$$? ;\
	exit $$tests_status_code ;\
