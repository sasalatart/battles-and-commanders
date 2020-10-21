compose_dev = docker-compose -f docker/compose-dev.yml -p batcoms_dev
compose_test = docker-compose -f docker/compose-test.yml -p batcoms_test
data_url = "https://bit.ly/3dOeqyZ"

.PHONY : help
help :
	@echo "help           : Runs this help command."
	@echo "build          : Builds the api, seeder and scraper bins targeting Linux."
	@echo "clean          : Removes the api, seeder and scraper bins created by build."
	@echo "dev_up         : [docker] turns on a Postgres database and the api in dev mode."
	@echo "dev_seed_local : [docker] runs the seeder for the dev api from local files."
	@echo "dev_seed_url   : [docker] runs the seeder for the dev api from a remote file (you can use the data_url option override)."
	@echo "dev_destroy    : [docker] stops and removes the dev containers."
	@echo "test_up        : [docker] turns on a Postgres database and the api in test mode."
	@echo "test_destroy   : [docker] stops and removes the test containers."
	@echo "test           : [docker] runs the test suites."
	@echo "scrape         : [docker] runs the scraper and stores results in data.json."

build:
	GOOS=linux go build -o api cmd/api/main.go
	GOOS=linux go build -o seeder cmd/seeder/main.go
	GOOS=linux go build -o scraper cmd/scraper/main.go

clean:
	rm api seeder scraper

dev_up:
	${compose_dev} up

dev_seed_local:
	${compose_dev} exec api go run cmd/seeder/main.go

dev_seed_url:
	${compose_dev} exec api go run cmd/seeder/main.go -dataURL=${data_url}

dev_destroy:
	${compose_dev} down && ${compose_dev} rm -f

test_up:
	${compose_test} up

test_destroy:
	${compose_test} down && ${compose_test} rm -f

test:
	${compose_test} exec api go test ./...

scrape:
	docker-compose -f docker/compose-dev-scraper.yml up
