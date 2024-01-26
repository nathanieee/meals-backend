DOCKER_COMPOSE := docker-compose
ENV_FILE := .env

ifndef ENV_EXISTS
$(shell if [ ! -f $(ENV_FILE) ]; then cp example.env $(ENV_FILE); fi)
endif

.PHONY: dev

dev:
	$(DOCKER_COMPOSE) down --volumes --remove-orphans
	$(DOCKER_COMPOSE) up --build