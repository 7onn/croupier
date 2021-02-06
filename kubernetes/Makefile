default: help

help:
	@echo 
	@echo "available commands"
	@echo "  - add-namespace"
	@echo "  - rm-namespace"
	@echo "  - db-port-forward #holds the terminal"
	@echo "  - db-connect #holds the terminal"
	@echo "  - add-metrics-server"
	@echo "  - rm-metrics-server"
	@echo "  - add-database"
	@echo "  - rm-database"
	@echo "  - add-redis"
	@echo "  - rm-redis"
	@echo "  - add-dsock-api"
	@echo "  - rm-dsock-api"
	@echo "  - add-dsock-worker"
	@echo "  - rm-dsock-worker"
	@echo "  - add-server"
	@echo "  - rm-server"
	@echo
	@echo "  - add-production (i.e: install production namespace and charts)"
	@echo "  - rm-production (i.e: remove production namespace and it's content)"
	@echo "  - print-logs (i.e: logs deploy/server-prd)"
	@echo
	@echo "note: you also might use 'add' scripts to update charts"

add-namespace:
	@kubectl create ns prd

rm-namespace:
	@kubectl delete ns prd

add-metrics-server:
	@kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

rm-metrics-server:
	@kubectl delete -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

add-database:
	@helm upgrade -i database-prd \
		--namespace prd \
		--values ./values/database/prd.yaml \
		./charts/database

rm-database:
	@helm -n prd del database-prd

db-port-forward:
	@kubectl -n prd port-forward svc/database-prd 5432:5432

db-connect:
	@psql "postgresql://postgres:dev@localhost:5432/dev"

add-redis:
	@helm upgrade -i redis-prd \
		--namespace prd \
		--values ./values/redis/prd.yaml \
		./charts/redis

rm-redis:
	@helm -n prd del redis-prd


add-dsock-api:
	@helm upgrade -i dsock-api-prd \
		--namespace prd \
		--values ./values/dsock-api/prd.yaml \
		./charts/dsock-api

rm-dsock-api:
	@helm -n prd del dsock-api-prd


add-dsock-worker:
	@helm upgrade -i dsock-worker-prd \
		--namespace prd \
		--values ./values/dsock-worker/prd.yaml \
		./charts/dsock-worker

rm-dsock-worker:
	@helm -n prd del dsock-worker-prd


add-server:
	@sleep 10
	@kubectl apply -f ./secrets/server/prd.secrets.yaml
	@helm upgrade -i server-prd \
		--namespace prd \
		--values ./values/server/prd.yaml \
		./charts/server

rm-server:
	@helm -n prd del server-prd
	@kubectl delete -f ./secrets/server/prd.secrets.yaml


.PHONY: add-production
add-production: add-namespace add-redis add-database add-dsock-api add-dsock-worker add-server
	@sleep 10
	@kubectl -n prd get pods

.PHONY: rm-production
rm-production: rm-server rm-dsock-worker rm-dsock-api rm-database rm-redis rm-namespace

print-logs:
	@kubectl -n prd logs deploy/server-prd
