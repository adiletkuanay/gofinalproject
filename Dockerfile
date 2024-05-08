FROM golang:1.20 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN go build -o /app/main ./cmd/canteen

FROM alpine:latest AS runtime

WORKDIR /app

COPY --from=build /app/main .

EXPOSE 8080

CMD ["./main"]
