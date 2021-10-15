.PHONY: all clean network start stop log login
all: login

clean:
	docker-compose down --volumes --rmi all

network:
	@if [ -z "`docker network ls | grep local-shared`" ]; then docker network create local-shared; fi

start: network
	docker-compose up -d --build

stop:
	docker-compose down

ps:
	docker-compose ps

log:
	docker-compose logs -f

login:
	mysql -h127.0.0.1 -uroot -Dmimosa --default-character-set=utf8
