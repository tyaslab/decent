FROM golang:1.25.0-alpine3.22 AS build-stage
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /server cmd/server/main.go

FROM alpine:3.22

WORKDIR /app
COPY --from=build-stage /server ./
EXPOSE 8080

CMD ["/app/server"]
