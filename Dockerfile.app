FROM golang:latest

RUN mkdir /app
WORKDIR /app
ENV CONFIG_PATH=./configs/remote.yaml
ENV POSTGRES_PASS=admin
ENV REDIS_PASS=admin

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main ./cmd/app/main.go

CMD [ "./main" ]