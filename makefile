clean:
	docker-compose -f docker-compose.yaml down -v

dev:
	if [ ! -f .env ]; then cp .env.example .env; fi;
	docker-compose -f docker-compose.yaml up --build

restart:
	docker-compose -f docker-compose.yaml down -v
	if [ ! -f .env ]; then cp .env.example .env; fi;	
	docker-compose -f docker-compose.yaml up --build