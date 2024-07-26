DOCKER_COMPOSE := docker-compose
DOCKER_COMPOSE_FILE := -f deployments/docker-compose.local.yaml
DOCKER_EXEC := docker exec
APP_CONTAINER := meals-go
ENV_FILE := .env

ifndef ENV_EXISTS
$(shell if [ ! -f $(ENV_FILE) ]; then cp example.env $(ENV_FILE); fi)
endif

.PHONY: dev

dev:
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_FILE) down --volumes --remove-orphans
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_FILE) up --build

generate-mockery:
	$(DOCKER_EXEC) -it $(APP_CONTAINER) /bin/sh -c "mockery --all"
