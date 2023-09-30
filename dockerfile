FROM golang:1.20-alpine
WORKDIR /app
COPY . .

RUN apk update
RUN apk add make

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN ln -s /go/bin/linux_amd64/migrate /usr/local/bin/migrate

CMD make migrate-up;go run main.go