FROM golang:1.22.4-bullseye as builder 

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify 

COPY . .

RUN test -d /storage || mkdir /storage

RUN CGO_ENABLED=1 go build -o ./main ./cmd/shortener/main.go

FROM debian:bullseye-slim  AS executer 

COPY --from=builder /app/main / 

EXPOSE 8082

CMD ["/main"]
