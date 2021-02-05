default: help

help:
	@echo 
	@echo "available commands"
	@echo "  - add-namespace"
	@echo "  - rm-namespace"
	@echo "  - add-metrics-server"
	@echo "  - rm-metrics-server"
	# @echo "  - add-database"
	# @echo "  - rm-database"
	@echo "  - add-redis"
	@echo "  - rm-redis"
	@echo "  - add-dsock-api"
	@echo "  - rm-dsock-api"
	# @echo "  - add-dsock-worker"
	# @echo "  - rm-dsock-worker"
	# @echo "  - add-server"
	# @echo "  - rm-server"
	# @echo "  - add-ingress (i.e: reverse proxy)" 
	# @echo "  - rm-ingress (i.e: reverse proxy)"
	# @echo
	# @echo "  - add-production (i.e: install production namespace and charts)"
	# @echo "  - rm-production (i.e: remove production namespace and it's content)"
	@echo
	# @echo "  - list-packages"
	# @echo "  - add-package [text=]"
	# @echo "  - rm-package [id=]"
	@echo
	# @echo "  - print-logs (print stdout and stderr from ingress container)"
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

# add-database:
# 	@helm upgrade -i database-prd \
# 		--namespace prd \
# 		--values ./values/database/prd.yaml \
# 		./charts/database

# rm-database:
# 	@helm -n prd del database-prd

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


# add-dsock-worker:
# 	@helm upgrade -i dsock-worker-prd \
# 		--namespace prd \
# 		--values ./values/dsock-worker/prd.yaml \
# 		./charts/dsock-worker

# rm-dsock-worker:
# 	@helm -n prd del dsock-worker-prd


# add-server:
# 	@helm upgrade -i server-prd \
# 		--namespace prd \
# 		--values ./values/server/prd.yaml \
# 		./charts/server

# rm-server:
# 	@helm -n prd del server-prd

# add-ingress:
# 	@helm upgrade -i ingress-prd \
# 		--namespace prd \
# 		--values ./values/ingress/prd.yaml \
# 		./charts/ingress

# rm-ingress:
# 	@helm -n prd del ingress-prd

# .PHONY: add-production
# add-production: add-namespace add-redis add-dsock-api add-database add-server add-ingress
# 	@sleep 10
# 	@kubectl -n prd get pods

# .PHONY: rm-production
# rm-production: rm-ingress rm-server rm-database add-dsock-api rm-namespace

# print-logs:
# 	@kubectl -n prd logs $(shell kubectl -n prd get pods -o name | grep ingress | sed -e 's/pod\///g' | awk '{print $$1}')

# list-packages:
# 	@curl -L minikube:32080/packages

# add-package:
# ifdef text
# 	@curl -L -X POST minikube:32080/packages -d $(text)
# else
# 	@echo 'parameter text is not defined'
# 	@echo 'e.g: make add-package text=DESCRIPTION'
# endif

# rm-package:
# ifdef id
# 	@curl -L -X DELETE minikube:32080/packages/$(id)
# else
# 	@echo 'parameter id not defined'
# 	@echo 'e.g: make rm-package id=ID'
# endif
