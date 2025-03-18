
go-mod:
	@go mod tidy
	@go mod vendor

up:
	@docker-compose up -d --force-recreate

down:
	@docker-compose down 
	@docker image rm optum-market-data-service --force
	@docker image rm optum-trader-service --force

restart:
	@make down && make up