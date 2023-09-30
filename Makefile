include .env

## Delete app
app-rm:
	docker rm -f app-container
	docker rmi app-image

## Build and run app
app-up:
	docker build -t app-image .
	docker run -dp 8080:8080\
		--network fio-service-app --network-alias app\
		-e POSTGRES_HOST=${POSTGRES_HOST}\
		-e POSTGRES_USER=${POSTGRES_USER}\
		-e POSTGRES_PASSWORD=${POSTGRES_PASSWORD}\
		-e POSTGRES_DB=${POSTGRES_DB}\
		-e POSTGRES_PORT=${POSTGRES_PORT}\
		-e REDIS_URL=${REDIS_URL}\
		-e KAFKA_CONSUMER_URL=${KAFKA_CONSUMER_URL}\
		-e KAFKA_PRODUCER_URL=${KAFKA_PRODUCER_URL}\
		-e NATIONALIZE_URL=${NATIONALIZE_URL}\
		-e AGIFY_URL=${AGIFY_URL}\
		-e GENDERIZE_URL=${GENDERIZE_URL}\
		--name app-container\
		app-image



## Run migrations UP
migrate-up:
	migrate -path ./migrations -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" -verbose up

## Run migrations DOWN
migrate-down:
	migrate -path ./migrations -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" -verbose down 1

## Create a DB migration files e.g `make migrate-create name=migration-name`
migrate-create:
	migrate create -ext sql  -dir ./migrations -seq $(name)