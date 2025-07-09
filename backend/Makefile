build:
	@go build -o bin/gobank
rebuild:	
	@rm bank.db
	@go build -o bin/gobank


run: build
	@./bin/gobank

docker-build:
	@docker build -t gobank .

docker-run:
	@docker run --rm -p 3000:3000 -e CLIENT="http://localhost:5500" gobank

