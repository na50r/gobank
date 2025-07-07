build:
	@go build -o bin/gobank
rebuild:	
	@rm bank.db
	@go build -o bin/gobank


run: build
	@./bin/gobank

seed:
	@./bin/gobank --seed

docker-build:
	@docker build -t gobank .

docker-run:
	@docker run --rm -p 3000:3000 gobank

