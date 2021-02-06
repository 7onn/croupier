default: help

help:
	@echo 
	@echo "available commands"
	@echo "  - minikube-start"
	@echo "  - minikube-stop"
	@echo "  - minikube-addons-enabled"
	@echo "  - minikube-addons-disabled"
	@echo

minikube-start:
	@minikube start
	@kubectl config use-context minikube

minikube-stop:
	@minikube stop

PHONY: minikube-addons-enabled
minikube-addons-enabled: minikube-addons-disabled
	@minikube addons enable ingress

minikube-addons-disabled:
	@minikube addons disable ingress
