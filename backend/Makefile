build:
	@go build -o bin/gobank
rebuild:	
	@rm store.db
	@go build -o bin/gobank

seed:
	@./bin/gobank --seed


run: build
	@./bin/gobank

docker-build:
	@docker build -t gobank .

docker-run:
	@docker run --rm -p 3000:3000 -e CLIENT="http://localhost:5500" gobank

