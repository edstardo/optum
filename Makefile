
go-mod:
	@go mod tidy
	@go mod vendor

up:
	@docker-compose up -d --force-recreate

down:
	@docker-compose down 
	@docker image rm mini-trader-market-data-service --force
	@docker image rm mini-trader-trader-service --force

restart:
	@make down && make up