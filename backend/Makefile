build:
	@go build -o bin/gobank
rebuild:	
	@rm bank.db
	@go build -o bin/gobank


run: build
	@./bin/gobank

seed:
	@./bin/gobank --seed
