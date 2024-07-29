DOCKER_COMPOSE := docker-compose
DOCKER_COMPOSE_FILE_LOCAL := -f deployments/docker-compose.local.yaml
DOCKER_COMPOSE_FILE_PROD := -f deployments/docker-compose.prod.yaml
DOCKER_EXEC := docker exec
APP_CONTAINER := meals-go
ENV_FILE := .env

ifndef ENV_EXISTS
$(shell if [ ! -f $(ENV_FILE) ]; then cp example.env $(ENV_FILE); fi)
endif

.PHONY: local

local:
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_FILE_LOCAL) --env-file ./.env down --volumes --remove-orphans
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_FILE_LOCAL) --env-file ./.env up --build

prod:
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_FILE_PROD) down --volumes --remove-orphans
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_FILE_PROD) up --build

generate-mockery:
	$(DOCKER_EXEC) -it $(APP_CONTAINER) /bin/sh -c "mockery --all"
