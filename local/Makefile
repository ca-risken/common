.PHONY: all clean start stop log login
all: login

clean:
	docker-compose down --volumes --rmi all

start:
	docker-compose up -d --build

stop:
	docker-compose down

ps:
	docker-compose ps

log:
	docker-compose logs -f

login:
	mysql -h127.0.0.1 -uroot -Dmimosa
