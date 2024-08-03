#---------- builder -------------------------
FROM golang:1.22.5-alpine
WORKDIR /var/www/html/service

RUN apk add build-base

COPY proto /var/www/html/proto

COPY ./service2/go.mod /var/www/html/service
RUN go mod download

COPY ./service2 /var/www/html/service

RUN cp configs.yaml.example configs.yaml

RUN CGO_ENABLED=0 go build -o main main.go

EXPOSE 50052

CMD ["./main"]

