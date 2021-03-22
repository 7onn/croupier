default: help

help:
	@echo 
	@echo "available commands"
	@echo "  - local-db"
	@echo "  - rm-db"
	@echo "  - db-deploy"
	@echo 

.PHONY: local-db
local-db:
	docker run -i --rm -d \
	--network host \
	--name db \
	-e POSTGRES_PASSWORD=dev \
	-e DATABASE_NAME=dev \
	devbytom/pg:latest

.PHONY: rm-db
rm-db:
	docker rm -f $$(docker ps -a -q --filter="name=db")

.PHONY: db-deploy
db-deploy:
	docker build -t devbytom/pg:latest .
	docker push devbytom/pg:latest
