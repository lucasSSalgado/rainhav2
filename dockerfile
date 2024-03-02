FROM golang:1.22-alpine AS base

WORKDIR /app
COPY . .

RUN go mod download && CGO_ENABLED=0 GOOS=linux go build -o /app/server

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=base /app/server /server

EXPOSE 8080
ENTRYPOINT ["./server"]